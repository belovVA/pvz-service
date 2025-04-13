package pgdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
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
	DB DB
}

func NewPVZRepository(db DB) *PVZRepository {
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
		return uuid.Nil, fmt.Errorf("%s: %w", FailedBuildQuery, err)
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
		return nil, fmt.Errorf("%s: %w", FailedBuildQuery, err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&pvz.ID,
		&pvz.RegistrationDate,
		&pvz.City,
	)
	if err != nil {
		return nil, fmt.Errorf(PvzNotFound)
	}

	return converter.ToPvzFromPvzRepo(&pvz), nil
}

func (r *PVZRepository) GetIDListPvz(ctx context.Context) ([]uuid.UUID, error) {
	query, args, err := sq.
		Select(pvzIDColumn).
		From(pvzTable).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", FailedBuildQuery, err)
	}
	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", FailedExecuteQuery, err)
	}

	defer rows.Close()
	result := make([]uuid.UUID, 0, 100)
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("%s: %w", FailedScanRow, err)
		}

		result = append(result, id)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%w", rows.Err())
	}

	return result, nil
}
