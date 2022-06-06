package book

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/infra/postgres"
	bookSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/book"
)

func DeleteBooks(ctx context.Context, bookIDs book.IDs, market market.Name, executioner postgres.QueryExecutioner) error {
	queryBuilder := squirrel.StatementBuilderType{}
	query, args, err := queryBuilder.
		Delete(bookSchema.TableName).
		Where(squirrel.Eq{
			bookSchema.ColumnID:     bookIDs,
			bookSchema.ColumnMarket: market,
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
