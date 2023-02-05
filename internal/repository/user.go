package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/acool-kaz/forum-auth-service/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func newUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (uint, error) {
	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, email, username, password) VALUES($1, $2, $3, $4, $5) RETURNING id;", userTable)

	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer prep.Close()

	var id uint

	if err = prep.QueryRowContext(ctx, user.FirstName, user.LastName, user.Email, user.Username, user.Password).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) Get(ctx context.Context) ([]models.User, error) {
	args := []interface{}{}
	argsStr := []string{}
	argsNum := 1

	email := ctx.Value(models.Email)
	if email != nil {
		argsStr = append(argsStr, fmt.Sprintf("email = $%d", argsNum))
		args = append(args, email.(string))
		argsNum++
	}

	username := ctx.Value(models.Username)
	if username != nil {
		argsStr = append(argsStr, fmt.Sprintf("username = $%d", argsNum))
		args = append(args, username.(string))
		argsNum++
	}

	whereCondition := ""
	if len(argsStr) != 0 {
		whereCondition = " WHERE " + strings.Join(argsStr, " AND ")
	}

	query := fmt.Sprintf("SELECT * FROM %s%s;", userTable, whereCondition)

	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	rows, err := prep.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
