package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"pvz-service/internal/repository/pgdb"
)

type Repository struct {
	*pgdb.UserRepository
	*pgdb.PVZRepository
	*pgdb.ReceptionRepository
	*pgdb.ProductRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		UserRepository:      pgdb.NewUserRepository(db),
		PVZRepository:       pgdb.NewPVZRepository(db),
		ReceptionRepository: pgdb.NewReceptionRepository(db),
		ProductRepository:   pgdb.NewProductRepository(db),
	}
}
