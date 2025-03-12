package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-jwt-hexagonal/internal/adapters/api"
	"go-jwt-hexagonal/internal/adapters/jwt"
	"go-jwt-hexagonal/internal/adapters/repository"
	"go-jwt-hexagonal/internal/core/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	// 1. Infrastructure/Adapters Layer
	userRepo := repository.NewInMemoryUserRepository()
	userService := services.NewUserService(userRepo)
	authService := jwt.NewJWTAdapter(jwtSecret, userService)
	apiHandler := api.NewAPIHandler(userService, authService)

	// 2. Wiring up the API (Primary Port - HTTP Handlers)
	router := mux.NewRouter()

	router.HandleFunc("/register", apiHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", apiHandler.LoginHandler).Methods("POST")

	// Protected endpoint - requires JWT authentication
	protectedRoute := router.PathPrefix("/protected").Subrouter()
	protectedRoute.Use(apiHandler.JwtAuthenticationMiddleware) // Apply middleware
	protectedRoute.HandleFunc("", apiHandler.ProtectedHandler).Methods("GET")

	// 3. Start the Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	fmt.Printf("Server listening on port :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
