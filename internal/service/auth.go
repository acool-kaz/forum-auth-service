package service

import (
	"context"
	"errors"

	"github.com/acool-kaz/forum-auth-service/internal/models"
	"github.com/acool-kaz/forum-auth-service/internal/repository"
)

type AuthService struct {
	userRepo repository.User
}

func newAuthService(userRepo repository.User) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, user models.User) (uint, error) {
	users, err := s.userRepo.Get(context.WithValue(ctx, models.Email, user.Email))
	if err != nil {
		return 0, err
	}
	if len(users) != 0 {
		return 0, errors.New("email exist")
	}

	_, err = s.userRepo.Get(context.WithValue(ctx, models.Username, user.Username))
	if err != nil {
		return 0, err
	}
	if len(users) != 0 {
		return 0, errors.New("username exist")
	}

	return s.userRepo.Create(ctx, user)
}
