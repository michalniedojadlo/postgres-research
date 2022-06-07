package book

import (
	"github.com/google/uuid"
)

type ID uuid.UUID
type IDs []ID

type ISBN string
type Author string

func (ids IDs) ToUUIDArray() []uuid.UUID {
	uuidArray := make([]uuid.UUID, 0, len(ids))

	for _, id := range ids {
		uuidArray = append(uuidArray, uuid.UUID(id))
	}

	return uuidArray
}
