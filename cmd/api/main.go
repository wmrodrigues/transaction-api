package main

import (
	"context"
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
	ctx := context.Background()
	db, err := config.GetPostgresConnection(ctx, &appConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	userRepository := postgres.NewUserRepository(db)

	// Services
	userService := service.NewUserService(userRepository)

	// Controllers
	userHandler := handler.NewUserHandler(userService)
	healthHandler := handler.NewHealthHandler()

	// Application HTTP routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("POST /users", userHandler.CreateUser)

	mux.HandleFunc("GET /health", healthHandler.GetHealth)

	// Starting the server
	log.Printf("Server started on port %s\n", appConfig.ServAddr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appConfig.ServAddr), mux))
}
