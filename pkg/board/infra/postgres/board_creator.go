package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/infra/postgres"

	boardSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/board"
)

type BoardCreator struct {
	client       *postgres.Client
	queryBuilder squirrel.StatementBuilderType
}

func NewBoardCreator(postgresClient *postgres.Client) *BoardCreator {
	return &BoardCreator{
		client:       postgresClient,
		queryBuilder: squirrel.StatementBuilderType{},
	}
}

func (boardCreator *BoardCreator) Create(ctx context.Context, id board.ID, market market.Name, name string) error {
	query, args, err := boardCreator.queryBuilder.
		Insert(boardSchema.TableName).
		Columns(
			boardSchema.ColumnID,
			boardSchema.ColumnMarket,
			boardSchema.ColumnName).
		Values(id,
			market,
			name).
		ToSql()
	if err != nil {
		return fmt.Errorf("creating query failed: %w", err)
	}

	conn, err := boardCreator.client.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection form the pool failed: %w", err)
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing query failed: %w", err)
	}

	return nil
}
