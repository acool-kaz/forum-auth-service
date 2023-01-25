package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/acool-kaz/forum-auth-service/internal/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	log.Println("init db")
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Name, cfg.Database.Password, cfg.Database.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	if err = createTables(db); err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	up, err := os.ReadFile("./migrations/up.sql")
	if err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	if _, err = db.Exec(string(up)); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	return nil
}
