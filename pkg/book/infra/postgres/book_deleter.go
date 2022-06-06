package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/infra/postgres"
	schema "github.com/brainly/postgres-research/internal/infra/postgres/schema/book"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
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
		Delete(schema.TableName).
		Where(squirrel.Eq{
			schema.ColumnID:     bookID,
			schema.ColumnMarket: market,
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
		Update(schema.TableName).
		Where(squirrel.Eq{
			schema.ColumnID: bookID,
		}).
		SetMap(map[string]interface{}{
			schema.ColumnDeletedAt: time.Now().UTC(),
			schema.ColumnDeletedBy: userID,
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
