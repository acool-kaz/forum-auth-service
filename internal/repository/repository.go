package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/acool-kaz/forum-auth-service/internal/models"
)

const (
	userTable = "users"
)

type User interface {
	Create(ctx context.Context, user models.User) (uint, error)
	Get(ctx context.Context) ([]models.User, error)
}

type Repository struct {
	User User
}

func InitRepository(db *sql.DB) *Repository {
	log.Println("init repository")

	return &Repository{
		User: newUserRepository(db),
	}
}
