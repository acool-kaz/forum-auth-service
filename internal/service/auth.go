package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/acool-kaz/forum-auth-service/internal/models"
	"github.com/acool-kaz/forum-auth-service/internal/repository"
	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	userRepo repository.User
}

func newAuthService(userRepo repository.User) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) ParseToken(ctx context.Context, accessToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.Token{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("JWT_SALT")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*models.Token)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.Id, nil
}

func (s *AuthService) Register(ctx context.Context, user models.User) (uint, error) {
	users, err := s.userRepo.Get(context.WithValue(ctx, models.Email, user.Email))
	if err != nil {
		return 0, err
	}
	if len(users) != 0 {
		return 0, errors.New("email exist")
	}

	users, err = s.userRepo.Get(context.WithValue(ctx, models.Username, user.Username))
	if err != nil {
		return 0, err
	}
	if len(users) != 0 {
		return 0, errors.New("username exist")
	}

	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, user models.User) (string, string, error) {
	var (
		curUser []models.User
		err     error
	)

	if user.Email != "" {
		curUser, err = s.userRepo.Get(context.WithValue(ctx, models.Email, user.Email))
	} else if user.Username != "" {
		curUser, err = s.userRepo.Get(context.WithValue(ctx, models.Username, user.Username))
	}

	if err != nil {
		return "", "", err
	}

	if len(curUser) == 0 || curUser[0].Password != user.Password {
		return "", "", errors.New("user not found")
	}

	access, err := newAccessToken(curUser[0].Id)
	if err != nil {
		return "", "", err
	}

	refresh, err := newRefreshToken(curUser[0].Id)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func newAccessToken(userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Token{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(models.AccessTokenTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id: userId,
	})

	return token.SignedString([]byte(os.Getenv("JWT_SALT")))
}

func newRefreshToken(userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Token{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(models.RefreshTokenTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id: userId,
	})

	return token.SignedString([]byte(os.Getenv("JWT_SALT")))
}
