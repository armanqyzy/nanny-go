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
	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключаемся к БД
	db, err := database.New(cfg.Database.ConnectionString())
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Инициализируем роутер
	r := mux.NewRouter()

	// Применяем middleware
	r.Use(middleware.CORS)

	// HTML страницы (ПЕРВЫМИ!)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	}).Methods("GET")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	}).Methods("GET")

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/register.html")
	}).Methods("GET")

	r.HandleFunc("/dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dashboard.html")
	}).Methods("GET")

	r.HandleFunc("/sitter-dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/sitter-dashboard.html")
	}).Methods("GET")

	r.HandleFunc("/admin-dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/admin-dashboard.html")
	}).Methods("GET")

	// Статические файлы (CSS, JS)
	staticDir := http.Dir("./static")
	staticHandler := http.FileServer(staticDir)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	// Инициализируем модули
	setupAuthModule(r, db)
	setupPetsModule(r, db)
	setupBookingsModule(r, db)
	setupReviewsModule(r, db)
	setupServicesModule(r, db)
	setupAdminModule(r, db)

	// Обратная совместимость со старыми URL
	authRepo := auth.NewRepository(db.DB)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)
	r.HandleFunc("/register/owner", authHandler.RegisterOwner).Methods("POST")
	r.HandleFunc("/register/sitter", authHandler.RegisterSitter).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Запускаем сервер
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("✅ Server running on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
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

	r.HandleFunc("/api/pets", handler.CreatePet).Methods("POST")
	r.HandleFunc("/api/pets/{id:[0-9]+}", handler.GetPet).Methods("GET")
	r.HandleFunc("/api/pets/{id:[0-9]+}", handler.UpdatePet).Methods("PUT")
	r.HandleFunc("/api/pets/{id:[0-9]+}", handler.DeletePet).Methods("DELETE")
	r.HandleFunc("/api/owners/{owner_id:[0-9]+}/pets", handler.GetOwnerPets).Methods("GET")
}

func setupBookingsModule(r *mux.Router, db *database.Database) {
	repo := bookings.NewRepository(db.DB)
	service := bookings.NewService(repo)
	handler := bookings.NewHandler(service)

	r.HandleFunc("/api/bookings", handler.CreateBooking).Methods("POST")
	r.HandleFunc("/api/bookings/{id:[0-9]+}", handler.GetBooking).Methods("GET")
	r.HandleFunc("/api/owners/{owner_id:[0-9]+}/bookings", handler.GetOwnerBookings).Methods("GET")
	r.HandleFunc("/api/sitters/{sitter_id:[0-9]+}/bookings", handler.GetSitterBookings).Methods("GET")
	r.HandleFunc("/api/bookings/{id:[0-9]+}/confirm", handler.ConfirmBooking).Methods("POST")
	r.HandleFunc("/api/bookings/{id:[0-9]+}/cancel", handler.CancelBooking).Methods("POST")
	r.HandleFunc("/api/bookings/{id:[0-9]+}/complete", handler.CompleteBooking).Methods("POST")
}

func setupReviewsModule(r *mux.Router, db *database.Database) {
	repo := reviews.NewRepository(db.DB)
	service := reviews.NewService(repo)
	handler := reviews.NewHandler(service)

	r.HandleFunc("/api/reviews", handler.CreateReview).Methods("POST")
	r.HandleFunc("/api/reviews/{id:[0-9]+}", handler.GetReview).Methods("GET")
	r.HandleFunc("/api/reviews/{id:[0-9]+}", handler.UpdateReview).Methods("PUT")
	r.HandleFunc("/api/reviews/{id:[0-9]+}", handler.DeleteReview).Methods("DELETE")
	r.HandleFunc("/api/sitters/{sitter_id:[0-9]+}/reviews", handler.GetSitterReviews).Methods("GET")
	r.HandleFunc("/api/sitters/{sitter_id:[0-9]+}/rating", handler.GetSitterRating).Methods("GET")
	r.HandleFunc("/api/bookings/{booking_id:[0-9]+}/review", handler.GetBookingReview).Methods("GET")
}

func setupServicesModule(r *mux.Router, db *database.Database) {
	repo := services.NewRepository(db.DB)
	service := services.NewService(repo)
	handler := services.NewHandler(service)

	r.HandleFunc("/api/services", handler.CreateService).Methods("POST")
	r.HandleFunc("/api/services/{id:[0-9]+}", handler.GetService).Methods("GET")
	r.HandleFunc("/api/services/{id:[0-9]+}", handler.UpdateService).Methods("PUT")
	r.HandleFunc("/api/services/{id:[0-9]+}", handler.DeleteService).Methods("DELETE")
	r.HandleFunc("/api/sitters/{sitter_id:[0-9]+}/services", handler.GetSitterServices).Methods("GET")
	r.HandleFunc("/api/services/search", handler.SearchServices).Methods("GET")
}

func setupAdminModule(r *mux.Router, db *database.Database) {
	repo := admin.NewRepository(db.DB)
	service := admin.NewService(repo)
	handler := admin.NewHandler(service)

	r.HandleFunc("/api/admin/sitters/pending", handler.GetPendingSitters).Methods("GET")
	r.HandleFunc("/api/admin/sitters/{sitter_id:[0-9]+}/approve", handler.ApproveSitter).Methods("POST")
	r.HandleFunc("/api/admin/sitters/{sitter_id:[0-9]+}/reject", handler.RejectSitter).Methods("POST")
	r.HandleFunc("/api/admin/sitters/{sitter_id:[0-9]+}", handler.GetSitterDetails).Methods("GET")
	r.HandleFunc("/api/admin/users", handler.GetAllUsers).Methods("GET")
	r.HandleFunc("/api/admin/users/{user_id:[0-9]+}", handler.GetUser).Methods("GET")
	r.HandleFunc("/api/admin/users/{user_id:[0-9]+}", handler.DeleteUser).Methods("DELETE")
}
