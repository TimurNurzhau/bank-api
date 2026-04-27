package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"bank-api/config"
	"bank-api/handlers"
	"bank-api/middleware"
	"bank-api/repositories"
	"bank-api/services"

	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Инициализация слоёв
	repos := repositories.NewRepositories(db)
	svcs := services.NewServices(repos, cfg)
	h := handlers.NewHandlers(svcs)

	// Шедулер
	scheduler := services.NewScheduler(repos, svcs.Email)
	scheduler.Start(12 * time.Hour)

	// Маршрутизатор
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.LoggingMiddleware)

	// Публичные эндпоинты
	r.HandleFunc("/register", h.Auth.Register).Methods("POST")
	r.HandleFunc("/login", h.Auth.Login).Methods("POST")

	// Защищённые эндпоинты
	authRouter := r.PathPrefix("").Subrouter()
	authRouter.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Счета
	authRouter.HandleFunc("/accounts", h.Account.Create).Methods("POST")
	authRouter.HandleFunc("/accounts", h.Account.List).Methods("GET")

	// Карты
	authRouter.HandleFunc("/cards", h.Card.Issue).Methods("POST")
	authRouter.HandleFunc("/cards", h.Card.List).Methods("GET")

	// Переводы
	authRouter.HandleFunc("/transfer", h.Transfer.Transfer).Methods("POST")
	authRouter.HandleFunc("/deposit", h.Transfer.Deposit).Methods("POST")

	// Кредиты
	authRouter.HandleFunc("/credits", h.Credit.Create).Methods("POST")
	authRouter.HandleFunc("/credits", h.Credit.List).Methods("GET")
	authRouter.HandleFunc("/credits/{creditId}/schedule", h.Credit.GetSchedule).Methods("GET")

	fmt.Printf("Server starting on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}