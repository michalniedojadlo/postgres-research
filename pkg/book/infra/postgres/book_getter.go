package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
	bookDomain "github.com/brainly/postgres-research/pkg/book"

	"github.com/brainly/postgres-research/internal/infra/postgres"
	bookSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/book"
	bookToBoardSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/booktoboard"
	bookToTopicSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/booktotopic"
)

type BookGetter struct {
	client       *postgres.Client
	queryBuilder squirrel.StatementBuilderType
}

func NewBookGetter(client *postgres.Client) *BookGetter {
	return &BookGetter{
		client:       client,
		queryBuilder: squirrel.StatementBuilderType{},
	}
}

func (bookGetter *BookGetter) GetBooks(ctx context.Context, market market.Name, topicID topic.ID, boardIDs board.IDs) ([]bookDomain.BookWithDetails, error) {
	query, args, err := bookGetter.queryBuilder.
		Select(
			bookSchema.ColumnID,
			bookSchema.ColumnAuthor,
			bookSchema.ColumnMarket,
			bookSchema.ColumnPublished,
			bookSchema.ColumnTitle,
			bookSchema.ColumnISBN,
			bookToBoardSchema.WithTableName(bookToBoardSchema.ColumnBoardID),
			bookToTopicSchema.WithTableName(bookToTopicSchema.ColumnTopicID),
		).
		From(bookSchema.TableName).
		Where(squirrel.Eq{
			bookSchema.ColumnMarket: market,
		}).
		LeftJoin("%s ON %s.%s = %s.%s", bookToTopicSchema.TableName, bookToTopicSchema.TableName, bookToTopicSchema.ColumnBookID, bookSchema.TableName, bookSchema.ColumnID).
		LeftJoin("%s ON %s.%s = %s.%s", bookToBoardSchema.TableName, bookToBoardSchema.TableName, bookToBoardSchema.ColumnBookID, bookSchema.TableName, bookSchema.ColumnID).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "building GetBooks query")
	}

	conn, err := bookGetter.client.GetConn(ctx)
	if err != nil {
		return nil, postgres.NewErrAcquiringConnection(err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "execution of GetBooks query")
	}

	bookList := make([]bookDomain.BookWithDetails, 0)
	for rows.Next() {
		var bookWithDetails bookSchema.BookWithDetails
		err = rows.Scan(
			&bookWithDetails.ID,
			&bookWithDetails.Author,
			&bookWithDetails.Market,
			&bookWithDetails.Published,
			&bookWithDetails.Title,
			&bookWithDetails.ISBN,
			&bookWithDetails.BoardIDs,
			&bookWithDetails.TopicID)
		if err != nil {
			return nil, errors.Wrap(err, "scanning result of GetBooks query")
		}

		bookList = append(bookList, ToDomainBookWithDetails(bookWithDetails))
	}

	return bookList, nil
}

func (bookGetter *BookGetter) GetBookByID(ctx context.Context, bookID book.ID) (*bookDomain.Book, error) {
	conn, err := bookGetter.client.GetConn(ctx)
	if err != nil {
		return nil, postgres.NewErrAcquiringConnection(err)
	}
	defer conn.Release()

	query, args, err := bookGetter.queryBuilder.
		Select(
			bookSchema.ColumnID,
			bookSchema.ColumnMarket,
			bookSchema.ColumnISBN,
			bookSchema.ColumnTitle,
			bookSchema.ColumnAuthor,
			bookSchema.ColumnPublished).
		From(bookSchema.TableName).
		Where(squirrel.Eq{
			bookSchema.ColumnID:        bookID,
			bookSchema.ColumnDeletedBy: nil,
			bookSchema.ColumnDeletedAt: nil,
		}).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building GetBookByID query")
	}

	book := bookSchema.Book{}
	err = conn.QueryRow(ctx, query, args...).
		Scan(&book.ID,
			&book.Market,
			&book.ISBN,
			&book.Title,
			&book.Author,
			&book.Published)
	if err != nil {
		return nil, fmt.Errorf("scanning query result failed: %w", err)
	}

	domainBook := ToDomainBook(book)
	return &domainBook, nil
}
