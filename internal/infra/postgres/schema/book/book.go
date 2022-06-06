package book

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID        uuid.UUID `db:"id"`
	Market    string    `db:"market"`
	Author    string    `db:"author"`
	ISBN      string    `db:"isbn"`
	Title     string    `db:"title"`
	Published bool      `db:"published"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
	DeletedAt time.Time `db:"deletedAt"`
	DeletedBy uuid.UUID `db:"deletedBy"`
}

const (
	TableName = `books`

	ColumnID        = `id`
	ColumnMarket    = `market`
	ColumnAuthor    = `author`
	ColumnISBN      = `isbn`
	ColumnTitle     = `title`
	ColumnPublished = `published`
	ColumnCreatedAt = `createdAt`
	ColumnUpdatedAt = `name`
	ColumnDeletedAt = `deletedAt`
	ColumnDeletedBy = `deletedBy`
)

type BookWithDetails struct {
	Book
	TopicID  uuid.UUID   `db:"topicId"`
	BoardIDs []uuid.UUID `db:"boardIds"`
}
