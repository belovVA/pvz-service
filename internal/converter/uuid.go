package converter

import "github.com/google/uuid"

func ParseUuid(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
