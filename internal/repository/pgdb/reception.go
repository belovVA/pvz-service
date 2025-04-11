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
	FailedCreateReception = "failed to Create Reception"
	ReceptionNotFound     = "reception not found"
)

const (
	receptionTable    = "reception"
	receptionIDColumn = "id"
	dateTimeColumn    = "date_time"
	isClosedStatus    = "is_closed"
	pvzIDColumnFK     = "pvz_id"
)

type ReceptionRepository struct {
	DB *pgxpool.Pool
}

func NewReceptionRepository(db *pgxpool.Pool) *ReceptionRepository {
	return &ReceptionRepository{
		DB: db,
	}
}

func (r *ReceptionRepository) CreateReception(ctx context.Context, pvzID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID

	query, args, err := sq.
		Insert(receptionTable).
		Columns(pvzIDColumnFK).
		Values(pvzID).
		Suffix("RETURNING " + receptionIDColumn).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %s", FailedCreateReception, err.Error())
	}

	return id, nil
}

func (r *ReceptionRepository) GetReceptionByID(ctx context.Context, id uuid.UUID) (*model.Reception, error) {
	var reception modelRepo.Reception

	query, args, err := sq.
		Select(receptionIDColumn, dateTimeColumn, isClosedStatus, pvzIDColumnFK).
		From(receptionTable).
		Where(sq.Eq{receptionIDColumn: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.IsClosedStatus,
		&reception.PvzID,
	)
	if err != nil {
		return nil, fmt.Errorf(ReceptionNotFound)
	}

	return converter.ToReceptionFromReceptionRepo(&reception), nil
}

func (r *ReceptionRepository) GetLastReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error) {
	var reception modelRepo.Reception

	query, args, err := sq.
		Select(receptionIDColumn, dateTimeColumn, isClosedStatus, pvzIDColumnFK).
		From(receptionTable).
		Where(sq.Eq{pvzIDColumnFK: pvzID}).
		OrderBy(dateTimeColumn + " DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.IsClosedStatus,
		&reception.PvzID,
	)
	if err != nil {
		return nil, fmt.Errorf(ReceptionNotFound)
	}

	return converter.ToReceptionFromReceptionRepo(&reception), nil
}
