package main

import (
	"fmt"
	"log"
	"net/http"
	config "transaction-api/internal/config"
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

	// Services
	userService := service.NewUserService(userRepository, transactionRepository)
	transactionService := service.NewTransactionService(transactionRepository)

	// Controllers
	userHandler := handler.NewUserHandler(userService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	healthHandler := handler.NewHealthHandler()

	// Application HTTP routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/users/{id}", userHandler.GetUser)
	mux.HandleFunc("POST /v1/users", userHandler.CreateUser)

	mux.HandleFunc("GET /v1/transactions/{id}", transactionHandler.GetTransaction)
	mux.HandleFunc("GET /v1/transactions", transactionHandler.GetTransactions)
	mux.HandleFunc("POST /v1/transactions", transactionHandler.CreateTransaction)

	mux.HandleFunc("GET /health", healthHandler.GetHealth)

	// Starting the server
	log.Printf("Server started on port %s\n", appConfig.ServAddr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appConfig.ServAddr), mux))
}
