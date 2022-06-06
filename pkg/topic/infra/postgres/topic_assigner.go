package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"

	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
	"github.com/brainly/postgres-research/internal/infra/postgres"

	internalTopic "github.com/brainly/postgres-research/internal/infra/postgres/topic"
)

type TopicAssigner struct {
	client       *postgres.Client
	queryBuilder squirrel.StatementBuilderType
}

func NewTopicAssigner(client *postgres.Client) *TopicAssigner {
	return &TopicAssigner{
		client:       client,
		queryBuilder: squirrel.StatementBuilderType{},
	}
}

func (topicAssigner *TopicAssigner) AssignToBook(ctx context.Context, topicID topic.ID, bookID book.ID, market market.Name) error {
	conn, err := topicAssigner.client.GetConn(ctx)
	if err != nil {
		return postgres.NewErrAcquiringConnection(err)
	}

	err = internalTopic.AssignTopicToBook(ctx, topicID, bookID, market, conn)
	if err != nil {
		return err
	}

	return nil
}
