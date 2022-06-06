package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
	"github.com/brainly/postgres-research/internal/infra/postgres"

	bookInfra "github.com/brainly/postgres-research/internal/infra/postgres/book"
	topicInfra "github.com/brainly/postgres-research/internal/infra/postgres/topic"
)

type TopicDeleter struct {
	client       *postgres.Client
	queryBuilder squirrel.StatementBuilderType
}

func NewTopicDeleter(client *postgres.Client) *TopicDeleter {
	return &TopicDeleter{
		client:       client,
		queryBuilder: squirrel.StatementBuilderType{},
	}
}

func (topicDeleter *TopicDeleter) DeleteTopic(ctx context.Context, topicID topic.ID, market market.Name) error {
	conn, err := topicDeleter.client.GetConn(ctx)
	if err != nil {
		return postgres.NewErrAcquiringConnection(err)
	}
	defer conn.Release()

	// // START Transaction
	tx, _ := conn.Begin(ctx)

	// find books to delete
	booksToDelete, err := bookInfra.GetBooks(ctx, nil, &topicID, nil, market, tx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// delete books
	bookIDs := make(book.IDs, 0, len(booksToDelete))
	for _, bookToDelete := range booksToDelete {
		bookIDs = append(bookIDs, bookToDelete.ID)
	}

	err = bookInfra.DeleteBooks(ctx, bookIDs, market, tx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// delete topic
	err = topicInfra.DeleteTopic(ctx, topicID, market, topicDeleter.queryBuilder, tx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// // END transaction

	err = tx.Commit(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("committing CreateBookWithAssignedFilters transaction: %w", err)
	}

	return nil
}
