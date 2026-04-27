package main

import (
	"fmt"
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
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	cfg := config.Load()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatalf("Failed to ping database: %v", err)
	}

	logger.Info("Successfully connected to PostgreSQL!")

	repos := repositories.NewRepositories(db)
	svcs := services.NewServices(repos, cfg, logger)
	h := handlers.NewHandlers(svcs, repos)

	scheduler := services.NewScheduler(repos, svcs.Email, logger)
	scheduler.Start(12 * time.Hour)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware(logger))

	r.HandleFunc("/register", h.Auth.Register).Methods("POST")
	r.HandleFunc("/login", h.Auth.Login).Methods("POST")

	authRouter := r.PathPrefix("").Subrouter()
	authRouter.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	authRouter.HandleFunc("/accounts", h.Account.Create).Methods("POST")
    authRouter.HandleFunc("/accounts", h.Account.List).Methods("GET")
	authRouter.HandleFunc("/cards", h.Card.Issue).Methods("POST")
	authRouter.HandleFunc("/cards", h.Card.List).Methods("GET")
	authRouter.HandleFunc("/cards/pay", h.Card.Pay).Methods("POST")  
	authRouter.HandleFunc("/transfer", h.Transfer.Transfer).Methods("POST")
	authRouter.HandleFunc("/deposit", h.Transfer.Deposit).Methods("POST")
	authRouter.HandleFunc("/credits", h.Credit.Create).Methods("POST")
	authRouter.HandleFunc("/credits", h.Credit.List).Methods("GET")
	authRouter.HandleFunc("/credits/{creditId}/schedule", h.Credit.GetSchedule).Methods("GET")
	authRouter.HandleFunc("/analytics", h.Analytics.GetAnalytics).Methods("GET")
	authRouter.HandleFunc("/accounts/{accountId}/predict", h.Analytics.PredictBalance).Methods("GET")

	logger.Infof("Server starting on port %s", cfg.ServerPort)
	logger.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}