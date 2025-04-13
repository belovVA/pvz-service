package pgdb_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pvz-service/internal/repository/pgdb"
)

func TestProductRepository_CreateProduct(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewProductRepository(mock)

	typeProduct := "Техника"
	receptionID := uuid.New()
	newID := uuid.New()

	mock.ExpectQuery(`INSERT INTO product\s*\(type_product,reception_id\)\s*VALUES\s*\(\$1,\$2\)\s*RETURNING id`).
		WithArgs(typeProduct, receptionID).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(newID))

	id, err := repo.CreateProduct(context.Background(), typeProduct, receptionID)
	require.NoError(t, err)
	assert.Equal(t, newID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_GetProductByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewProductRepository(mock)

	productID := uuid.New()
	receptionID := uuid.New()
	now := time.Now().UTC()
	typeProduct := "Одежда"

	mock.ExpectQuery(`SELECT id, date_time, type_product, reception_id FROM product WHERE id = \$1`).
		WithArgs(productID.String()).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "date_time", "type_product", "reception_id"}).
				AddRow(productID, now, typeProduct, receptionID),
		)

	product, err := repo.GetProductByID(context.Background(), productID)
	require.NoError(t, err)
	assert.Equal(t, productID, product.ID)
	assert.WithinDuration(t, now, product.DateTime, time.Second)
	assert.Equal(t, typeProduct, product.TypeProduct)
	assert.Equal(t, receptionID, product.ReceptionID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_GetLastProduct(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewProductRepository(mock)

	receptionID := uuid.New()
	productID := uuid.New()
	now := time.Now().UTC()
	typeProduct := "Книги"

	mock.ExpectQuery(`SELECT id, date_time, type_product, reception_id FROM product WHERE reception_id = \$1 ORDER BY date_time DESC LIMIT 1`).
		WithArgs(receptionID.String()).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "date_time", "type_product", "reception_id"}).
				AddRow(productID, now, typeProduct, receptionID),
		)

	product, err := repo.GetLastProduct(context.Background(), receptionID)
	require.NoError(t, err)
	assert.Equal(t, productID, product.ID)
	assert.Equal(t, typeProduct, product.TypeProduct)
	assert.Equal(t, receptionID, product.ReceptionID)
	assert.WithinDuration(t, now, product.DateTime, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_DeleteProductByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewProductRepository(mock)

	productID := uuid.New()

	mock.ExpectExec(`DELETE FROM product WHERE id = \$1`).
		WithArgs(productID.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.DeleteProductByID(context.Background(), productID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_GetProductSliceByReceptionID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewProductRepository(mock)

	receptionID := uuid.New()
	productID := uuid.New()
	now := time.Now().UTC()
	typeProduct := "Техника"

	mock.ExpectQuery(`SELECT id, date_time, type_product, reception_id FROM product WHERE reception_id = \$1`).
		WithArgs(receptionID.String()).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "date_time", "type_product", "reception_id"}).
				AddRow(productID, now, typeProduct, receptionID),
		)

	products, err := repo.GetProductSliceByReceptionID(context.Background(), receptionID)
	require.NoError(t, err)
	require.Len(t, products, 1)
	assert.Equal(t, productID, products[0].ID)
	assert.Equal(t, typeProduct, products[0].TypeProduct)
	assert.WithinDuration(t, now, products[0].DateTime, time.Second)
	assert.Equal(t, receptionID, products[0].ReceptionID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
