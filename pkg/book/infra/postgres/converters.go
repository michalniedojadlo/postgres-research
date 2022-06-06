package postgres

import (
	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"

	bookSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/book"
	bookDomain "github.com/brainly/postgres-research/pkg/book"
)

func ToDomainBook(bookModel bookSchema.Book) bookDomain.Book {
	return bookDomain.Book{
		ID:        book.ID(bookModel.ID),
		Market:    market.Name(bookModel.Market),
		Author:    book.Author(bookModel.Author),
		ISBN:      book.ISBN(bookModel.ISBN),
		Title:     bookModel.Title,
		Published: bookModel.Published,
	}
}

func ToDomainBookWithDetails(bookWithDetailsModel bookSchema.BookWithDetails) bookDomain.BookWithDetails {
	boardIDs := make(board.IDs, 0, len(bookWithDetailsModel.BoardIDs))
	for _, boardID := range bookWithDetailsModel.BoardIDs {
		boardIDs = append(boardIDs, board.ID(boardID))
	}

	return bookDomain.BookWithDetails{
		Book:   ToDomainBook(bookWithDetailsModel.Book),
		Topic:  topic.ID(bookWithDetailsModel.TopicID),
		Boards: boardIDs,
	}
}
