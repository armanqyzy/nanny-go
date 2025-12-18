package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"nanny-backend/internal/admin"
	"nanny-backend/internal/auth"
	"nanny-backend/internal/bookings"
	"nanny-backend/internal/common/database"
	"nanny-backend/internal/common/middleware"
	"nanny-backend/internal/pets"
	"nanny-backend/internal/reviews"
	"nanny-backend/internal/services"
	"nanny-backend/pkg/config"
)

func main() {
	cfg := config.Load()

	db, err := connectWithRetry(cfg.Database.ConnectionString(), 10, 3*time.Second)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	setupAuthModule(r, db)
	setupPetsModule(r, db)
	setupBookingsModule(r, db)
	setupReviewsModule(r, db)
	setupServicesModule(r, db)
	setupAdminModule(r, db)

	frontendDir := "../nanny-front"
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(frontendDir)))

	handler := middleware.CORS(
		middleware.RequestLogger(
			middleware.RateLimit(r),
		),
	)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		startBookingExpirationWorker(ctx, db)
	}()

	go func() {
		log.Printf("✅ API server started on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down application...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ Forced shutdown: %v", err)
	} else {
		log.Println("✅ HTTP server stopped gracefully")
	}

	wg.Wait()
	log.Println("✅ Background worker stopped, application exited cleanly")
}

func connectWithRetry(dsn string, attempts int, delay time.Duration) (*database.Database, error) {
	var db *database.Database
	var err error

	for i := 1; i <= attempts; i++ {
		db, err = database.New(dsn)
		if err == nil {
			log.Println("✅ Connected to database")
			return db, nil
		}

		log.Printf("⏳ DB connection failed (attempt %d/%d): %v", i, attempts, err)
		time.Sleep(delay)
	}

	return nil, err
}

func startBookingExpirationWorker(ctx context.Context, db *database.Database) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	log.Println("🔄 Background worker started: booking expiration checker")

	for {
		select {
		case <-ctx.Done():
			log.Println("⏹️ Background worker stopped")
			return

		case <-ticker.C:
			checkExpiredBookings(ctx, db)
		}
	}
}

func checkExpiredBookings(ctx context.Context, db *database.Database) {
	query := `
		UPDATE bookings
		SET status = 'cancelled'
		WHERE status = 'pending'
		  AND start_time < NOW() - INTERVAL '24 hours'
	`

	res, err := db.DB.ExecContext(ctx, query)
	if err != nil {
		log.Printf("❌ Worker error updating expired bookings: %v", err)
		return
	}

	affected, _ := res.RowsAffected()
	if affected > 0 {
		log.Printf("✅ Worker cancelled %d expired booking(s)", affected)
	}
}

func setupAuthModule(r *mux.Router, db *database.Database) {
	repo := auth.NewRepository(db.DB)
	service := auth.NewService(repo)
	handler := auth.NewHandler(service)

	r.HandleFunc("/api/auth/register/owner", handler.RegisterOwner).Methods("POST")
	r.HandleFunc("/api/auth/register/sitter", handler.RegisterSitter).Methods("POST")
	r.HandleFunc("/api/auth/login", handler.Login).Methods("POST")
}

func setupPetsModule(r *mux.Router, db *database.Database) {
	repo := pets.NewRepository(db.DB)
	service := pets.NewService(repo)
	handler := pets.NewHandler(service)

	r.Handle("/api/pets",
		middleware.AuthMiddleware(http.HandlerFunc(handler.CreatePet)),
	).Methods("POST")

	r.Handle("/api/pets/{id:[0-9]+}",
		middleware.AuthMiddleware(http.HandlerFunc(handler.UpdatePet)),
	).Methods("PUT")

	r.Handle("/api/pets/{id:[0-9]+}",
		middleware.AuthMiddleware(http.HandlerFunc(handler.DeletePet)),
	).Methods("DELETE")

	r.HandleFunc("/api/pets/{id:[0-9]+}", handler.GetPet).Methods("GET")
	r.HandleFunc("/api/owners/{owner_id:[0-9]+}/pets", handler.GetOwnerPets).Methods("GET")
}

func setupBookingsModule(r *mux.Router, db *database.Database) {
	repo := bookings.NewRepository(db.DB)
	service := bookings.NewService(repo)
	handler := bookings.NewHandler(service)

	r.Handle("/api/bookings",
		middleware.AuthMiddleware(http.HandlerFunc(handler.CreateBooking)),
	).Methods("POST")

	r.Handle("/api/bookings/{id:[0-9]+}/confirm",
		middleware.AuthMiddleware(http.HandlerFunc(handler.ConfirmBooking)),
	).Methods("POST")

	r.Handle("/api/bookings/{id:[0-9]+}/cancel",
		middleware.AuthMiddleware(http.HandlerFunc(handler.CancelBooking)),
	).Methods("POST")

	r.Handle("/api/bookings/{id:[0-9]+}/complete",
		middleware.AuthMiddleware(http.HandlerFunc(handler.CompleteBooking)),
	).Methods("POST")

	r.HandleFunc("/api/bookings/{id:[0-9]+}", handler.GetBooking).Methods("GET")
	r.HandleFunc("/api/owners/{owner_id:[0-9]+}/bookings", handler.GetOwnerBookings).Methods("GET")
	r.HandleFunc("/api/sitters/{sitter_id:[0-9]+}/bookings", handler.GetSitterBookings).Methods("GET")
}

func setupReviewsModule(r *mux.Router, db *database.Database) {
	repo := reviews.NewRepository(db.DB)
	service := reviews.NewService(repo)
	handler := reviews.NewHandler(service)

	r.Handle("/api/reviews",
		middleware.AuthMiddleware(http.HandlerFunc(handler.CreateReview)),
	).Methods("POST")

	r.Handle("/api/reviews/{id:[0-9]+}",
		middleware.AuthMiddleware(http.HandlerFunc(handler.UpdateReview)),
	).Methods("PUT")

	r.Handle("/api/reviews/{id:[0-9]+}",
		middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteReview)),
	).Methods("DELETE")

	r.HandleFunc("/api/reviews/{id:[0-9]+}", handler.GetReview).Methods("GET")
	r.HandleFunc("/api/sitters/{sitter_id:[0-9]+}/reviews", handler.GetSitterReviews).Methods("GET")
	r.HandleFunc("/api/sitters/{sitter_id:[0-9]+}/rating", handler.GetSitterRating).Methods("GET")
	r.HandleFunc("/api/bookings/{booking_id:[0-9]+}/review", handler.GetBookingReview).Methods("GET")
}

func setupServicesModule(r *mux.Router, db *database.Database) {
	repo := services.NewRepository(db.DB)
	service := services.NewService(repo)
	handler := services.NewHandler(service)

	r.HandleFunc("/api/services/search", handler.SearchServices).Methods("GET")
	r.HandleFunc("/api/sitters/{sitter_id:[0-9]+}/services", handler.GetSitterServices).Methods("GET")
	r.HandleFunc("/api/services/{id:[0-9]+}", handler.GetService).Methods("GET")

	r.Handle("/api/services",
		middleware.AuthMiddleware(http.HandlerFunc(handler.CreateService)),
	).Methods("POST")

	r.Handle("/api/services/{id:[0-9]+}",
		middleware.AuthMiddleware(http.HandlerFunc(handler.UpdateService)),
	).Methods("PUT")

	r.Handle("/api/services/{id:[0-9]+}",
		middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteService)),
	).Methods("DELETE")
}

func setupAdminModule(r *mux.Router, db *database.Database) {
	repo := admin.NewRepository(db.DB)
	service := admin.NewService(repo)
	handler := admin.NewHandler(service)

	r.Handle("/api/admin/sitters/pending",
		middleware.AuthMiddleware(http.HandlerFunc(handler.GetPendingSitters)),
	).Methods("GET")

	r.Handle("/api/admin/sitters/{sitter_id:[0-9]+}/approve",
		middleware.AuthMiddleware(http.HandlerFunc(handler.ApproveSitter)),
	).Methods("POST")

	r.Handle("/api/admin/sitters/{sitter_id:[0-9]+}/reject",
		middleware.AuthMiddleware(http.HandlerFunc(handler.RejectSitter)),
	).Methods("POST")

	r.Handle("/api/admin/sitters/{sitter_id:[0-9]+}",
		middleware.AuthMiddleware(http.HandlerFunc(handler.GetSitterDetails)),
	).Methods("GET")

	r.Handle("/api/admin/users",
		middleware.AuthMiddleware(http.HandlerFunc(handler.GetAllUsers)),
	).Methods("GET")

	r.Handle("/api/admin/users/{user_id:[0-9]+}",
		middleware.AuthMiddleware(http.HandlerFunc(handler.GetUser)),
	).Methods("GET")

	r.Handle("/api/admin/users/{user_id:[0-9]+}",
		middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteUser)),
	).Methods("DELETE")
}
