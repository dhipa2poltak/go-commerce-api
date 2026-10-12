package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/dhipa2poltak/go-commerce-api/internal/delivery/http/dto"
	"github.com/dhipa2poltak/go-commerce-api/internal/domain"
	"github.com/dhipa2poltak/go-commerce-api/internal/infrastructure/auth"
	"github.com/dhipa2poltak/go-commerce-api/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthUsecase struct {
	userRepository repository.UserRepository
	jwtService     *auth.JWTService
}

func NewAuthUsecase(
	userRepository repository.UserRepository,
	jwtService *auth.JWTService,
) *AuthUsecase {
	return &AuthUsecase{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, request dto.RegisterRequest) (*dto.AuthResponse, error) {
	name := strings.TrimSpace(request.Name)
	email := strings.ToLower(strings.TrimSpace(request.Email))

	if name == "" || email == "" || request.Password == "" {
		return nil, ErrInvalidInput
	}

	if len(request.Password) < 8 {
		return nil, ErrInvalidInput
	}

	_, err := u.userRepository.FindByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := auth.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken: accessToken,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (u *AuthUsecase) Login(ctx context.Context, request dto.LoginRequest) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))

	if email == "" || request.Password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !auth.CheckPassword(request.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken: accessToken,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
