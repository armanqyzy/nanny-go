package main

import (
	"fmt"
	"log"
	"net/http"

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

	db, err := database.New(cfg.Database.ConnectionString())
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	setupAuthModule(r, db)
	setupPetsModule(r, db)
	setupBookingsModule(r, db)
	setupReviewsModule(r, db)
	setupServicesModule(r, db)
	setupAdminModule(r, db)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("✅ API server running on http://localhost%s\n", addr)
	handler := middleware.CORS(r)
	log.Fatal(http.ListenAndServe(addr, handler))
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
