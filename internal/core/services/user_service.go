package services

import (
	"context"
	"errors"
	"go-jwt-hexagonal/internal/core/domain"
	"go-jwt-hexagonal/internal/core/ports"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo ports.UserRepository
}

// NewUserService creates a new UserService with the given UserRepository.
func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{userRepo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, email, password string) (*domain.User, error) {
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *userService) LoginUser(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}
