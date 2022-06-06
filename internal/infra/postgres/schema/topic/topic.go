package topic

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Market    string    `db:"market"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}

const (
	TableName = `books`

	ColumnID        = `id`
	ColumnName      = `name`
	ColumnMarket    = `market`
	ColumnCreatedAt = `createdAt`
	ColumnUpdatedAt = `name`
)
