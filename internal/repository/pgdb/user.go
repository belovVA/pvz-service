package pgdb

import (
	"context"
	"fmt"

	"pvz-service/internal/model"
	"pvz-service/internal/repository/pgdb/converter"
	modelRepo "pvz-service/internal/repository/pgdb/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const (
	FailedCreateUser = "failed to Create User"
	UserNotFound     = "user not found"
)

const (
	usersTable     = "users"
	userIDColumn   = "id"
	emailColumn    = "email"
	passwordColumn = "password"
	roleColumn     = "role"
)

type UserRepository struct {
	DB DB
}

func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	var id uuid.UUID

	query, args, err := sq.
		Insert(usersTable).
		Columns(emailColumn, passwordColumn, roleColumn).
		Values(user.Email, user.Password, user.Role).
		Suffix("RETURNING " + userIDColumn).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", FailedBuildQuery, err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %s", FailedCreateUser, err.Error())
	}

	return id, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user modelRepo.User

	query, args, err := sq.
		Select(userIDColumn, emailColumn, passwordColumn, roleColumn).
		From(usersTable).
		Where(sq.Eq{emailColumn: email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", FailedBuildQuery, err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", UserNotFound, email)
	}

	return converter.ToUserFromUserRepo(&user), nil
}
