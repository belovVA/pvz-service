package pgdb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"pvz-service/internal/model"
	"pvz-service/internal/repository/pgdb/converter"
	modelRepo "pvz-service/internal/repository/pgdb/model"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (uuid.UUID, error) {
	var id uuid.UUID

	query := `
        INSERT INTO users (email, password, role)
        VALUES ($1, $2, $3)
        RETURNING id;
    `
	// QueryRow выполняет запрос и ожидает одну строку в результате.
	err := r.DB.QueryRow(ctx, query, user.Email, user.Password, user.Role).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to Create User: %w", err)
	}

	return id, nil
}

func (r *UserRepository) Get(ctx context.Context, email, password string) (*model.User, error) {
	var user modelRepo.User
	query := `
		SELECT id, email, password, role
		FROM users
		WHERE email = $1 AND password = $2;
	`
	err := r.DB.QueryRow(ctx, query, email, password).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("не найден пользователь с указанными данными: %w", err)
	}
	return converter.ToUserFromUserRepo(&user), nil
}
