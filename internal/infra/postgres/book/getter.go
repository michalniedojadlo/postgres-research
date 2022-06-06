package book

import (
	"context"

	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
	"github.com/brainly/postgres-research/internal/infra/postgres"

	bookDomain "github.com/brainly/postgres-research/pkg/book"
)

func GetBooks(ctx context.Context, bookIDs book.IDs, topicID *topic.ID, boardIDs board.IDs, market market.Name, executioner postgres.QueryExecutioner) ([]bookDomain.Book, error) {
	return nil, nil
}
