package repository

import (
	"context"
	"database/sql"
	"errors"
	"maxwellzp/blog-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO user (username, email, password) VALUES (?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	user.ID, err = result.LastInsertId()
	return err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, username, email, password FROM user WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)

	user := &model.User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
