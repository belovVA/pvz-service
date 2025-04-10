package pgdb

import (
	"context"
	"fmt"

	"pvz-service/internal/model"
	"pvz-service/internal/repository/pgdb/converter"
	modelRepo "pvz-service/internal/repository/pgdb/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	FailedCreateUser = "failed to Create User"
	UserNotFound     = "user not found"
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
		return uuid.Nil, fmt.Errorf("%s: %s", FailedCreateUser, err.Error())
	}

	return id, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user modelRepo.User
	query := `
		SELECT id, email, password, role
		FROM users
		WHERE email = $1;
	`
	err := r.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", UserNotFound, email)
	}
	return converter.ToUserFromUserRepo(&user), nil
}
