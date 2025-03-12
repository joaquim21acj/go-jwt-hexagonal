package jwt

import (
	"go-jwt-hexagonal/internal/core/ports"
	"go-jwt-hexagonal/internal/core/services"
)

// NewJWTAdapter creates a new AuthService implementation using JWT.
func NewJWTAdapter(secretKey string, userService ports.UserService) ports.AuthService {
	return services.NewAuthService(secretKey, userService)
}
