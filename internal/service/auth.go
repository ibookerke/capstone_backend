package service

import (
	"context"
	"fmt"

	"github.com/ibookerke/capstone_backend/internal/config"
	"github.com/ibookerke/capstone_backend/internal/domain"
	"github.com/ibookerke/capstone_backend/internal/pkg/trm"
)

type AuthService struct {
	trm            trm.Manager
	userRepository domain.UserRepository
	cfg            config.Auth
}

type AuthResult struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func NewAuthService(
	trm trm.Manager,
	userRepository domain.UserRepository,
	cfg config.Auth,
) *AuthService {
	return &AuthService{
		trm:            trm,
		userRepository: userRepository,
		cfg:            cfg,
	}
}

func (a *AuthService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := a.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (a *AuthService) Register(ctx context.Context, email, password, name, surname string) (*domain.User, error) {
	user := domain.User{
		Email:    email,
		Password: password,
		Name:     name,
		Surname:  surname,
	}

	existingUser, err := a.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("register user error: %w", err)
	}
	if existingUser.ID != 0 {
		return nil, fmt.Errorf("user with provided email already exists")
	}

	_, err = a.userRepository.CreateUser(ctx, &user)

	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}

	return &user, nil
}
