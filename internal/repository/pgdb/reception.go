package pgdb

import (
	"context"
	"fmt"
	"time"

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
		return uuid.Nil, fmt.Errorf("%s: %w", BuildingQueryFailed, err)
	}

	if err = r.DB.QueryRow(ctx, query, args...).Scan(&id); err != nil {
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
		return nil, fmt.Errorf("%s: %w", BuildingQueryFailed, err)
	}

	if err = r.DB.QueryRow(ctx, query, args...).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.IsClosedStatus,
		&reception.PvzID,
	); err != nil {
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
		return nil, fmt.Errorf("%s: %w", BuildingQueryFailed, err)
	}

	if err = r.DB.QueryRow(ctx, query, args...).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.IsClosedStatus,
		&reception.PvzID,
	); err != nil {
		return nil, fmt.Errorf(ReceptionNotFound)
	}

	return converter.ToReceptionFromReceptionRepo(&reception), nil
}

func (r *ReceptionRepository) CloseReception(ctx context.Context, receptionID uuid.UUID) error {
	query, args, err := sq.
		Update(receptionTable).
		Set(isClosedStatus, true).
		Where(sq.Eq{receptionIDColumn: receptionID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос
	cmdTag, err := r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	// Проверяем, что была затронута хотя бы одна строка
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected, reception not found")
	}
	return nil
}

func (r *ReceptionRepository) GetReceptionsSliceWithTimeRange(ctx context.Context, begin time.Time, end time.Time) ([]model.Reception, error) {
	var result []model.Reception

	queryBuilder := sq.
		Select(receptionIDColumn, dateTimeColumn, isClosedStatus, pvzIDColumnFK).
		From(receptionTable).
		PlaceholderFormat(sq.Dollar)

	if !begin.IsZero() && !end.IsZero() {
		queryBuilder = queryBuilder.Where(sq.And{
			sq.GtOrEq{dateTimeColumn: begin},
			sq.LtOrEq{dateTimeColumn: end},
		})
	} else if !begin.IsZero() {
		queryBuilder = queryBuilder.Where(sq.GtOrEq{dateTimeColumn: begin})
	} else if !end.IsZero() {
		queryBuilder = queryBuilder.Where(sq.LtOrEq{dateTimeColumn: end})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", BuildingQueryFailed, err)
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var receptionRepo modelRepo.Reception
		if err = rows.Scan(
			&receptionRepo.ID,
			&receptionRepo.DateTime,
			&receptionRepo.IsClosedStatus,
			&receptionRepo.PvzID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		reception := converter.ToReceptionFromReceptionRepo(&receptionRepo)
		result = append(result, *reception)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", rows.Err())
	}

	return result, nil
}
