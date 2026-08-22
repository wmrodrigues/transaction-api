package main

import (
	"fmt"
	"log"
	"net/http"
	"transaction-api/internal/authentication"
	config "transaction-api/internal/config"
	"transaction-api/internal/database"
	"transaction-api/internal/handler"
	"transaction-api/internal/repository/postgres"
	"transaction-api/internal/service"
)

func main() {
	appConfig := config.LoadConfigs()
	db, err := config.GetDatabaseConnection(&appConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	userRepository := postgres.NewUserRepository(db)
	transactionRepository := postgres.NewTransactionRepository(db)

	// Transaction manager
	transactionManager := database.NewGormManager(db)

	// Services
	userService := service.NewUserService(userRepository, transactionRepository, transactionManager)
	transactionService := service.NewTransactionService(transactionRepository, userRepository, transactionManager)
	jwtService := authentication.NewJWTService(appConfig.JwtSecret)

	// Controllers
	userHandler := handler.NewUserHandler(userService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(userService, jwtService)

	// Authentication middleware
	authMiddleware := authentication.AuthMiddleware(jwtService)

	// Application HTTP routes
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/auth/tokens", authHandler.Tokens)

	mux.Handle("GET /v1/users/me", authMiddleware(http.HandlerFunc(userHandler.Me)))
	mux.Handle("GET /v1/users/{id}", authMiddleware(http.HandlerFunc(userHandler.GetUser)))
	mux.HandleFunc("POST /v1/users", userHandler.CreateUser)
	mux.Handle("GET /v1/users/{id}/transactions", authMiddleware(http.HandlerFunc(transactionHandler.GetTransactions)))
	mux.Handle("GET /v1/users/{id}/balance", authMiddleware(http.HandlerFunc(transactionHandler.GetUserBalance)))

	mux.Handle("GET /v1/transactions/{id}", authMiddleware(http.HandlerFunc(transactionHandler.GetTransaction)))
	mux.Handle("POST /v1/transactions", authMiddleware(http.HandlerFunc(transactionHandler.CreateTransaction)))

	mux.HandleFunc("GET /health", healthHandler.GetHealth)

	// Starting the server
	log.Printf("Server started on port %s\n", appConfig.ServAddr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appConfig.ServAddr), mux))
}
