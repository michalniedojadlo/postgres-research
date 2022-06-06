package book

import (
	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
	"github.com/google/uuid"
)

type CreateBookCMD struct {
	Market    market.Name
	ISBN      book.ISBN
	Author    book.Author
	Published bool
	Title     string
	UserID    uuid.UUID
}

type CreateBookWithDetailsCMD struct {
	CreateBookCMD
	Topic  topic.ID
	Boards board.IDs
	UserID uuid.UUID
}

type DeleteBookCMD struct {
	ID     book.ID
	Market market.Name
	UserID uuid.UUID
}

type SoftDeleteBookCMD struct {
	ID     book.ID
	Market market.Name
	UserID uuid.UUID
}
