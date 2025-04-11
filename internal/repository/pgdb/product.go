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
	FailedCreateProduct = "failed to Create Product"
	productNotFound     = "product not found"
)

const (
	productTable          = "product"
	productIDColumn       = "id"
	dateTimeProductColumn = "date_time"
	typeProductColumn     = "type_product"
	receptionIDFKColumn   = "reception_id"
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		DB: db,
	}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, typeProduct string, recepID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID

	query, args, err := sq.
		Insert(productTable).
		Columns(typeProductColumn, receptionIDFKColumn).
		Values(typeProduct, recepID).
		Suffix("RETURNING " + productIDColumn).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %s", FailedCreateProduct, err.Error())
	}
	return id, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	var product modelRepo.Product

	query, args, err := sq.
		Select(productIDColumn, dateTimeProductColumn, typeProductColumn, receptionIDFKColumn).
		From(productTable).
		Where(sq.Eq{productIDColumn: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&product.ID,
		&product.DateTime,
		&product.TypeProduct,
		&product.ReceptionID,
	)
	if err != nil {
		return nil, fmt.Errorf(productNotFound)
	}

	return converter.ToProductFromProductRepo(&product), nil
}

func (r *ProductRepository) GetLastProduct(ctx context.Context, receptionID uuid.UUID) (*model.Product, error) {
	var product modelRepo.Product

	query, args, err := sq.
		Select(productIDColumn, dateTimeProductColumn, typeProductColumn, receptionIDFKColumn).
		From(productTable).
		Where(sq.Eq{receptionIDFKColumn: receptionID}).
		OrderBy(dateTimeColumn + " DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&product.ID,
		&product.DateTime,
		&product.TypeProduct,
		&product.ReceptionID,
	)
	if err != nil {
		return nil, fmt.Errorf(productNotFound)
	}

	return converter.ToProductFromProductRepo(&product), nil
}

func (r *ProductRepository) DeleteProductByID(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Delete(productTable).
		Where(sq.Eq{productIDColumn: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	// Выполняем запрос
	cmdTag, err := r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute delete: %w", err)
	}

	// Проверяем, что была затронута хотя бы одна строка
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected, product not found")
	}

	return nil
}
