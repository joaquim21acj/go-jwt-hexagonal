package services

import (
	"context"
	"go-jwt-hexagonal/internal/core/domain"
	"go-jwt-hexagonal/internal/core/ports"
	"reflect"
	"testing"
)

func TestNewAuthService(t *testing.T) {
	type args struct {
		secretKey   string
		userService ports.UserService
	}
	tests := []struct {
		name string
		args args
		want ports.AuthService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthService(tt.args.secretKey, tt.args.userService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authService_GenerateToken(t *testing.T) {
	type fields struct {
		jwtSecretKey string
		userService  ports.UserService
	}
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authService{
				jwtSecretKey: tt.fields.jwtSecretKey,
				userService:  tt.fields.userService,
			}
			got, err := s.GenerateToken(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authService_ValidateToken(t *testing.T) {
	type fields struct {
		jwtSecretKey string
		userService  ports.UserService
	}
	type args struct {
		ctx         context.Context
		tokenString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authService{
				jwtSecretKey: tt.fields.jwtSecretKey,
				userService:  tt.fields.userService,
			}
			got, err := s.ValidateToken(tt.args.ctx, tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
