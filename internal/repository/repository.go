package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"pvz-service/internal/repository/pgdb"
)

type Repository struct {
	*pgdb.UserRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		UserRepository: pgdb.NewUserRepository(db),
	}
}
