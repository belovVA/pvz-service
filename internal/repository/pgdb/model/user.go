package modelRepo

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Role     string    `db:"role"`
}
