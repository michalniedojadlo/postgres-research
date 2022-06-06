package topic

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
	"github.com/brainly/postgres-research/internal/infra/postgres"

	topicSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/booktotopic"
)

func AssignTopicToBook(ctx context.Context, topicID topic.ID, bookID book.ID, market market.Name, executioner postgres.QueryExecutioner) error {
	queryBuilder := squirrel.StatementBuilderType{}
	query, args, err := queryBuilder.
		Insert(topicSchema.TableName).
		Columns(
			topicSchema.ColumnBookID,
			topicSchema.ColumnTopicID,
			topicSchema.ColumnMarket).
		Values(
			bookID,
			topicID,
			market).ToSql()
	if err != nil {
		return fmt.Errorf("building query failed: %w", err)
	}

	_, err = executioner.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	return nil
}

func ClearTopicToBookAssignments(ctx context.Context, bookID book.ID, market market.Name, executioner postgres.QueryExecutioner) error {
	queryBuilder := squirrel.StatementBuilderType{}
	query, args, err := queryBuilder.
		Delete(topicSchema.TableName).
		Where(squirrel.Eq{
			topicSchema.ColumnBookID: bookID,
			topicSchema.ColumnMarket: market,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("creating query failed: %w", err)
	}

	_, err = executioner.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing query failed: %w", err)
	}

	return nil
}
