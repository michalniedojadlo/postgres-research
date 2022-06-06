package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/infra/postgres"
	bookSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/book"
)

type ClearBoardsAssignmentFunc func(ctx context.Context, bookID book.ID, market market.Name, executioner postgres.QueryExecutioner) error
type ClearTopicAssignmentFunc func(ctx context.Context, bookID book.ID, market market.Name, executioner postgres.QueryExecutioner) error

type BookDeleter struct {
	client       *postgres.Client
	queryBuilder squirrel.StatementBuilderType

	clearBoardsAssignments ClearBoardsAssignmentFunc
	clearTopicAssignment   ClearTopicAssignmentFunc
}

func NewBookDeleter(client *postgres.Client, clearBoardsAssignmentsFunc ClearBoardsAssignmentFunc, clearTopicAssignment ClearTopicAssignmentFunc) *BookDeleter {
	return &BookDeleter{
		client:                 client,
		queryBuilder:           squirrel.StatementBuilderType{},
		clearBoardsAssignments: clearBoardsAssignmentsFunc,
		clearTopicAssignment:   clearTopicAssignment,
	}
}

func (bookDeleter *BookDeleter) DeleteBookByID(ctx context.Context, bookID book.ID, market market.Name) error {
	query, args, err := bookDeleter.queryBuilder.
		Delete(bookSchema.TableName).
		Where(squirrel.Eq{
			bookSchema.ColumnID:     bookID,
			bookSchema.ColumnMarket: market,
		}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building DeleteBookByID query")
	}

	conn, err := bookDeleter.client.GetConn(ctx)
	if err != nil {
		return postgres.NewErrAcquiringConnection(err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction failed: %w", err)
	}

	if err = bookDeleter.clearTopicAssignment(ctx, bookID, market, tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err = bookDeleter.clearBoardsAssignments(ctx, bookID, market, tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("executing query failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction failed: %w", err)
	}

	return nil
}

func (bookDeleter *BookDeleter) SoftDelete(ctx context.Context, bookID book.ID, userID uuid.UUID) error {
	query, args, err := bookDeleter.queryBuilder.
		Update(bookSchema.TableName).
		Where(squirrel.Eq{
			bookSchema.ColumnID: bookID,
		}).
		SetMap(map[string]interface{}{
			bookSchema.ColumnDeletedAt: time.Now().UTC(),
			bookSchema.ColumnDeletedBy: userID,
		}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building SoftDelete query")
	}

	conn, err := bookDeleter.client.GetConn(ctx)
	if err != nil {
		return postgres.NewErrAcquiringConnection(err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "execution of SoftDelete query")
	}

	return nil
}
