package book

import (
	"context"

	"github.com/google/uuid"

	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
)

type Book struct {
	ID        book.ID
	Market    market.Name
	Author    book.Author
	ISBN      book.ISBN
	Title     string
	Published bool
}

type BookWithDetails struct {
	Book
	Topic  topic.ID
	Boards board.IDs
}

type Service struct {
	bookCreator Creator
	bookDeleter Deleter
	bookGetter  Getter
}

func NewService(bookCreator Creator, bookDeleter Deleter, bookGetter Getter) *Service {
	return &Service{
		bookCreator: bookCreator,
		bookDeleter: bookDeleter,
		bookGetter:  bookGetter,
	}
}

type Creator interface {
	CreateBookAndAssignFilters(ctx context.Context, bookWithDetails *BookWithDetails) error
	CreateBook(ctx context.Context, book *Book) error
}

type Deleter interface {
	DeleteBookByID(ctx context.Context, bookID book.ID, market market.Name) error
	SoftDelete(ctx context.Context, bookID book.ID, userID uuid.UUID) error
}

type Getter interface {
	GetBookByID(ctx context.Context, bookID book.ID) (*Book, error)
	GetBooks(ctx context.Context, market market.Name, topicID topic.ID, boardIDs board.IDs) ([]BookWithDetails, error)
}
