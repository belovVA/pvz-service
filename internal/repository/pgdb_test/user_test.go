package pgdb_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pvz-service/internal/model"
	"pvz-service/internal/repository/pgdb"
)

func TestUserRepository_CreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewUserRepository(mock)

	t.Run("успешное создание пользователя", func(t *testing.T) {
		expectedID := uuid.New()
		user := &model.User{
			Email:    "test@example.com",
			Password: "hashedpass",
			Role:     "admin",
		}

		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(user.Email, user.Password, user.Role).
			WillReturnRows(mock.NewRows([]string{"id"}).AddRow(expectedID))

		id, err := repo.CreateUser(context.Background(), user)
		require.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("ошибка при создании пользователя", func(t *testing.T) {
		user := &model.User{
			Email:    "fail@example.com",
			Password: "failpass",
			Role:     "user",
		}

		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(user.Email, user.Password, user.Role).
			WillReturnError(errors.New("insert error"))

		id, err := repo.CreateUser(context.Background(), user)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
	})
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := pgdb.NewUserRepository(mock)

	t.Run("успешное получение пользователя", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       uuid.New(),
			Email:    "find@example.com",
			Password: "securepass",
			Role:     "admin",
		}

		mock.ExpectQuery(`SELECT id, email, password, role FROM users`).
			WithArgs(expectedUser.Email).
			WillReturnRows(mock.NewRows([]string{"id", "email", "password", "role"}).
				AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password, expectedUser.Role))

		user, err := repo.GetUserByEmail(context.Background(), expectedUser.Email)
		require.NoError(t, err)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, expectedUser.Role, user.Role)
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		email := "missing@example.com"

		mock.ExpectQuery(`SELECT id, email, password, role FROM users`).
			WithArgs(email).
			WillReturnError(errors.New("no rows in result set"))

		user, err := repo.GetUserByEmail(context.Background(), email)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
