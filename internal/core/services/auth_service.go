package services

import (
	"context"
	"errors"
	"fmt"
	"go-jwt-hexagonal/internal/core/domain"
	"go-jwt-hexagonal/internal/core/ports"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	jwtSecretKey string
	userService  ports.UserService
}

// NewAuthService creates a new AuthService with the JWT secret key.
func NewAuthService(secretKey string, userService ports.UserService) ports.AuthService {
	return &authService{jwtSecretKey: secretKey, userService: userService}
}

func (s *authService) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":    time.Now().Unix(),                     // Issued At Time
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}
	userID := int(userIDFloat)

	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found for token validation: %w", err)
	}

	return user, nil
}
