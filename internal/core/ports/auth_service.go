package ports

import (
	"context"
	"go-jwt-hexagonal/internal/core/domain"
)

// AuthService defines the operations for authentication-related business logic.
type AuthService interface {
	GenerateToken(ctx context.Context, user *domain.User) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (*domain.User, error)
}
