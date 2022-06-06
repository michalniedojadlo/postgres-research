package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
	"github.com/brainly/postgres-research/internal/infra/postgres"
	bookSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/book"

	bookDomain "github.com/brainly/postgres-research/pkg/book"
)

type BoardAssignerFunc func(ctx context.Context, boardIDs board.IDs, bookID book.ID, market market.Name, executioner postgres.QueryExecutioner) error
type TopicAssignerFunc func(ctx context.Context, topicID topic.ID, bookID book.ID, market market.Name, executioner postgres.QueryExecutioner) error

type BookCreator struct {
	client       *postgres.Client
	queryBuilder squirrel.StatementBuilderType

	assignBoards BoardAssignerFunc
	assignTopic  TopicAssignerFunc
}

func NewBookCreator(client *postgres.Client, boardAssigner BoardAssignerFunc, topicAssigner TopicAssignerFunc) *BookCreator {
	return &BookCreator{
		client:       client,
		queryBuilder: squirrel.StatementBuilderType{},
		assignBoards: boardAssigner,
		assignTopic:  topicAssigner,
	}
}

func (bookCreator *BookCreator) CreateBook(ctx context.Context, book *bookDomain.Book) error {
	conn, err := bookCreator.client.GetConn(ctx)
	if err != nil {
		return postgres.NewErrAcquiringConnection(err)
	}
	defer conn.Release()

	return bookCreator.createBook(ctx, book, conn)
}

func (bookCreator *BookCreator) CreateBookAndAssignFilters(ctx context.Context, book *bookDomain.BookWithDetails) error {
	conn, err := bookCreator.client.GetConn(ctx)
	if err != nil {
		return postgres.NewErrAcquiringConnection(err)
	}
	defer conn.Release()

	tx, _ := conn.Begin(ctx)
	if err = bookCreator.createBook(ctx, &book.Book, tx); err != nil {
		return errors.Wrap(err, "creating book")
	}

	err = bookCreator.assignTopic(ctx, book.Topic, book.ID, book.Market, tx)
	if err != nil {
		tx.Rollback(ctx)
		return errors.Wrap(err, "assigning book to topics")
	}

	err = bookCreator.assignBoards(ctx, book.Boards, book.ID, book.Market, tx)
	if err != nil {
		tx.Rollback(ctx)
		return errors.Wrap(err, "assigning book to boards")
	}

	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return errors.Wrap(err, "committing CreateBookWithAssignedFilters transaction")
	}

	return nil
}

func (bookCreator *BookCreator) createBook(ctx context.Context, book *bookDomain.Book, queryExecutioner postgres.QueryExecutioner) error {
	query, args, err := bookCreator.queryBuilder.
		Insert(bookSchema.TableName).
		Columns(
			bookSchema.ColumnID,
			bookSchema.ColumnMarket,
			bookSchema.ColumnISBN,
			bookSchema.ColumnTitle,
			bookSchema.ColumnAuthor,
			bookSchema.ColumnPublished).
		Values(
			book.ID,
			book.Market,
			book.ISBN,
			book.Title,
			book.Author,
			book.Published).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building CreateBook query")
	}

	row := queryExecutioner.QueryRow(ctx, query, args...)
	err = row.Scan(&book.ID)
	if err != nil {
		return errors.Wrap(err, "scanning result of CreateBook query")
	}

	return nil
}
