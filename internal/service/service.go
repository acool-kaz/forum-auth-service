package service

import (
	"context"
	"log"

	"github.com/acool-kaz/forum-auth-service/internal/models"
	"github.com/acool-kaz/forum-auth-service/internal/repository"
)

type Auth interface {
	Register(ctx context.Context, user models.User) (uint, error)
}

type Service struct {
	Auth Auth
}

func InitService(repo *repository.Repository) *Service {
	log.Println("init service")
	return &Service{
		Auth: newAuthService(repo.User),
	}
}
