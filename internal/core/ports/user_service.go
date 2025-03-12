package ports

import (
	"context"
	"go-jwt-hexagonal/internal/core/domain"
)

// UserService defines the operations for user-related business logic.
type UserService interface {
	RegisterUser(ctx context.Context, email, password string) (*domain.User, error)
	LoginUser(ctx context.Context, email, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error) // For UserRepository
}
