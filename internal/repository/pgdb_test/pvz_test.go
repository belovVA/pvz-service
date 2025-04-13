package pgdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pvz-service/internal/repository/pgdb"
)

func TestPVZRepository_CreatePvz(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewPVZRepository(mock)

	expectedID := uuid.New()
	city := "Москва"

	mock.ExpectQuery("INSERT INTO pvz").
		WithArgs(city).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(expectedID))

	id, err := repo.CreatePvz(context.Background(), city)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPVZRepository_GetPvzByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewPVZRepository(mock)

	id := uuid.New()

	// Мокаем SQL-запрос на получение с ошибкой (если запись не найдена)
	mock.ExpectQuery(`^SELECT id, registration_date, city FROM pvz WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(fmt.Errorf("pvz not found"))

	// Пытаемся получить запись по несуществующему ID
	pvz, err := repo.GetPvzByID(context.Background(), id)
	assert.Error(t, err)                       // Ожидаем ошибку
	assert.Nil(t, pvz)                         // Ожидаем, что возвращаемый объект nil
	assert.EqualError(t, err, "pvz not found") // Проверяем, что ошибка именно такая, как мы ожидаем

}

func TestPVZRepository_GetPvzByID_AfterCreate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewPVZRepository(mock)

	id := uuid.New()
	city := "Казань"
	registrationDate := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Мокаем SQL-запрос на добавление нового PVZ
	mock.ExpectQuery(`INSERT INTO pvz \(city\) VALUES \(\$1\) RETURNING id`).
		WithArgs(city).
		WillReturnRows(
			pgxmock.NewRows([]string{"id"}).
				AddRow(id),
		)

	createdID, err := repo.CreatePvz(context.Background(), city)
	assert.NoError(t, err)
	assert.Equal(t, id, createdID)

	mock.ExpectQuery(`SELECT id, registration_date, city FROM pvz WHERE id = \$1`).
		WithArgs(id.String()).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "registration_date", "city"}).
				AddRow(id, registrationDate, city),
		)

	pvz, err := repo.GetPvzByID(context.Background(), createdID)
	assert.NoError(t, err)
	require.NotNil(t, pvz)
	assert.Equal(t, id, pvz.ID)
	assert.WithinDuration(t, registrationDate, pvz.RegistrationDate, time.Second)
	assert.Equal(t, city, pvz.City)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPVZRepository_GetIDListPvz(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewPVZRepository(mock)

	id1 := uuid.New()
	id2 := uuid.New()

	mock.ExpectQuery("SELECT id FROM pvz").
		WillReturnRows(
			pgxmock.NewRows([]string{"id"}).
				AddRow(id1).
				AddRow(id2),
		)

	ids, err := repo.GetIDListPvz(context.Background())
	assert.NoError(t, err)
	assert.Len(t, ids, 2)
	assert.Contains(t, ids, id1)
	assert.Contains(t, ids, id2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
