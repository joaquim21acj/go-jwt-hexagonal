package repository

import (
	"context"
	"errors"
	"go-jwt-hexagonal/internal/core/domain"
	"go-jwt-hexagonal/internal/core/ports"
	"sync"
)

type inMemoryUserRepository struct {
	users  map[int]*domain.User
	mu     sync.RWMutex // For thread safety if needed
	nextID int
}

// NewInMemoryUserRepository creates a new in-memory UserRepository.
func NewInMemoryUserRepository() ports.UserRepository {
	return &inMemoryUserRepository{
		users:  make(map[int]*domain.User),
		nextID: 1,
	}
}

func (r *inMemoryUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.users[r.nextID] = user
	r.nextID++
	return user, nil
}

func (r *inMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *inMemoryUserRepository) FindByID(ctx context.Context, id int) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
