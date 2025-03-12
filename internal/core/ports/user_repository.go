package ports

import (
	"context"
	"go-jwt-hexagonal/internal/core/domain"
)

// UserRepository defines the operations for user data persistence.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
}
