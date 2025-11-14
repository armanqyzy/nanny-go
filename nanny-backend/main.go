package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	var err error
	connStr := "postgres://postgres:Ana4aBada$$@localhost:5432/nanny_db?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("❌ Не удалось подключиться к БД:", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatal("❌ Таблица users не найдена:", err)
	}

	fmt.Println("✅ Подключено к БД. В таблице users:", count, "записей")

	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	staticDir := http.Dir("./static")
	staticHandler := http.FileServer(staticDir)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	r.HandleFunc("/register/owner", registerOwner).Methods("POST")
	r.HandleFunc("/register/sitter", registerSitter).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	}).Methods("GET")

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/register.html")
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	}).Methods("GET")

	fmt.Println("✅ Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func registerOwner(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "неверные данные"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	_, err := db.Exec(`INSERT INTO users (full_name, email, phone, password_hash, role)
                       VALUES ($1, $2, $3, $4, 'owner')`,
		input.FullName, input.Email, input.Phone, string(hashed))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "владелец зарегистрирован успешно"})
}

func registerSitter(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName        string `json:"full_name"`
		Email           string `json:"email"`
		Phone           string `json:"phone"`
		Password        string `json:"password"`
		ExperienceYears int    `json:"experience_years"`
		Certificates    string `json:"certificates"`
		Preferences     string `json:"preferences"`
		Location        string `json:"location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "неверные данные"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	var sitterID int
	err := db.QueryRow(`
        INSERT INTO users (full_name, email, phone, password_hash, role)
        VALUES ($1, $2, $3, $4, 'sitter')
        RETURNING user_id
    `, input.FullName, input.Email, input.Phone, string(hashed)).Scan(&sitterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, err = db.Exec(`
        INSERT INTO sitters (sitter_id, experience_years, certificates, preferences, location, status)
        VALUES ($1, $2, $3, $4, $5, 'pending')
    `, sitterID, input.ExperienceYears, input.Certificates, input.Preferences, input.Location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "няня зарегистрирована, ожидает подтверждения"})
}

func login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "неверные данные"})
		return
	}

	var id int
	var role string
	var hash string

	err := db.QueryRow(`SELECT user_id, role, password_hash FROM users WHERE email=$1`, input.Email).
		Scan(&id, &role, &hash)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "пользователь не найден"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "неверный пароль"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "вход выполнен",
		"role":    role,
	})
}
