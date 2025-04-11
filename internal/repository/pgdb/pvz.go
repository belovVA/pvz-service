package pgdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"pvz-service/internal/model"
	"pvz-service/internal/repository/pgdb/converter"
	modelRepo "pvz-service/internal/repository/pgdb/model"
)

const (
	FailedCreatePvz = "failed to Create Pvz"
	PvzNotFound     = "pvz not found"
)

const (
	pvzTable               = "pvz"
	pvzIDColumn            = "id"
	dateRegistrationColumn = "registration_date"
	cityColumn             = "city"
)

type PVZRepository struct {
	DB *pgxpool.Pool
}

func NewPVZRepository(db *pgxpool.Pool) *PVZRepository {
	return &PVZRepository{
		DB: db,
	}
}

func (r *PVZRepository) CreatePvz(ctx context.Context, city string) (uuid.UUID, error) {
	var id uuid.UUID

	query, args, err := sq.
		Insert(pvzTable).
		Columns(cityColumn).
		Values(city).
		Suffix("RETURNING " + pvzIDColumn).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %s", FailedCreatePvz, err.Error())
	}

	return id, nil
}

func (r *PVZRepository) GetPvzByID(ctx context.Context, id uuid.UUID) (*model.Pvz, error) {
	var pvz modelRepo.Pvz

	query, args, err := sq.
		Select(pvzIDColumn, dateRegistrationColumn, cityColumn).
		From(pvzTable).
		Where(sq.Eq{pvzIDColumn: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&pvz.ID,
		&pvz.RegistrationDate,
		&pvz.City,
	)
	if err != nil {
		return nil, fmt.Errorf("%s", PvzNotFound)
	}

	return converter.ToPvzFromPvzRepo(&pvz), nil
}
