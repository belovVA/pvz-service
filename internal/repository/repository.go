package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"pvz-service/internal/model"
	"pvz-service/internal/repository/pgdb"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (uuid.UUID, error)
	Get(ctx context.Context, email, password string) (*model.User, error)
}

type Repository struct {
	UserRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		UserRepository: pgdb.NewUserRepository(db),
	}
}
