package api

import (
	"context"
	"encoding/json"
	"go-jwt-hexagonal/internal/core/domain"
	"go-jwt-hexagonal/internal/core/ports"
	"net/http"
)

type apiHandler struct {
	userService ports.UserService
	authService ports.AuthService
}

// NewAPIHandler creates a new API handler with UserService and AuthService.
func NewAPIHandler(userService ports.UserService, authService ports.AuthService) *apiHandler {
	return &apiHandler{userService: userService, authService: authService}
}

// RegisterHandler handles user registration.
func (h *apiHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user, err := h.userService.RegisterUser(context.Background(), request.Email, request.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"userID":  user.ID,
		"email":   user.Email,
	})
}

// LoginHandler handles user login and JWT generation.
func (h *apiHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user, err := h.userService.LoginUser(context.Background(), request.Email, request.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := h.authService.GenerateToken(context.Background(), user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token":  token,
		"userID": user.ID,
		"email":  user.Email,
	})
}

// ProtectedHandler is a sample protected endpoint (requires JWT).
func (h *apiHandler) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*domain.User)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Protected endpoint accessed!",
		"userID":  user.ID,
		"email":   user.Email,
	})
}

// JwtAuthenticationMiddleware is middleware to validate JWT token.
func (h *apiHandler) JwtAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing JWT token")
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		user, err := h.authService.ValidateToken(context.Background(), tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid JWT token: "+err.Error())
			return
		}

		// Token is valid, add user to context for downstream handlers
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper functions for responses
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
