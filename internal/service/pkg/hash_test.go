package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Run("successfully hashes password", func(t *testing.T) {
		password := "123qweASD-1d"

		hashed, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)

		err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
		assert.NoError(t, err)
	})
}
