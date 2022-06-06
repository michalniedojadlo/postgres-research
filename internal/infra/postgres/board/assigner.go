package board

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/infra/postgres"

	boardSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/booktoboard"
)

func AssignBoardsToBook(ctx context.Context, boardIDs board.IDs, bookID book.ID, market market.Name, executioner postgres.QueryExecutioner) error {
	values := make([]interface{}, 0, len(boardIDs)*2)
	for _, boardID := range boardIDs {
		values = append(values, bookID)
		values = append(values, boardID)
		values = append(values, market)
	}

	queryBuilder := squirrel.StatementBuilderType{}

	query, args, err := queryBuilder.
		Insert(boardSchema.TableName).
		Columns(
			boardSchema.ColumnBookID,
			boardSchema.ColumnBoardID,
			boardSchema.ColumnMarket).
		Values(values...).ToSql()
	if err != nil {
		return fmt.Errorf("building query failed: %w", err)
	}

	_, err = executioner.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	return nil
}

func ClearBoardToBookAssignments(ctx context.Context, bookID book.ID, market market.Name, executioner postgres.QueryExecutioner) error {
	queryBuilder := squirrel.StatementBuilderType{}
	query, args, err := queryBuilder.
		Delete(boardSchema.TableName).
		Where(squirrel.Eq{
			boardSchema.ColumnBookID: bookID,
			boardSchema.ColumnMarket: market,
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
