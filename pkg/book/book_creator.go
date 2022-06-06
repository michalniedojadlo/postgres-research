package book

import (
	"context"
	"github.com/google/uuid"

	"github.com/brainly/postgres-research/internal/core/book"
)

type CreatorService struct {
}

func (service *Service) CreateBook(ctx context.Context, cmd *CreateBookCMD) error {
	newBookID := uuid.New()
	book := &Book{
		ID:        book.ID(newBookID),
		Market:    cmd.Market,
		Author:    cmd.Author,
		ISBN:      cmd.ISBN,
		Title:     cmd.Title,
		Published: cmd.Published,
	}

	err := service.bookCreator.CreateBook(ctx, book)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) CreateBookWithDetails(ctx context.Context, cmd *CreateBookWithDetailsCMD) error {
	newBookID := uuid.New()

	bookWithDetails := &BookWithDetails{
		Book: Book{
			ID:        book.ID(newBookID),
			Market:    cmd.Market,
			Author:    cmd.Author,
			ISBN:      cmd.ISBN,
			Title:     cmd.Title,
			Published: cmd.Published,
		},
		Topic:  cmd.Topic,
		Boards: cmd.Boards,
	}

	err := service.bookCreator.CreateBookAndAssignFilters(ctx, bookWithDetails)
	if err != nil {
		return err
	}

	return nil
}
