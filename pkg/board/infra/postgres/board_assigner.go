package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/infra/postgres"
	commonBoard "github.com/brainly/postgres-research/internal/infra/postgres/board"
)

type BoardAssigner struct {
	client       *postgres.Client
	queryBuilder squirrel.StatementBuilderType
}

func NewBoardAssigner(postgresClient *postgres.Client) *BoardAssigner {
	return &BoardAssigner{
		client:       postgresClient,
		queryBuilder: squirrel.StatementBuilderType{},
	}
}

func (boardAssigner *BoardAssigner) AssignToBook(ctx context.Context, boardIDs board.IDs, bookID book.ID, market market.Name) error {
	conn, err := boardAssigner.client.GetConn(ctx)
	if err != nil {
		return postgres.NewErrAcquiringConnection(err)
	}

	err = commonBoard.AssignBoardsToBook(ctx, boardIDs, bookID, market, conn)
	if err != nil {
		return err
	}

	return nil
}
