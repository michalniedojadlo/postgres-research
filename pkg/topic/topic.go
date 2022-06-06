package topic

import (
	"context"

	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
)

type Topic struct {
	ID     topic.ID
	Name   string
	Market market.Name
}

type Service struct {
	deleter  Deleter
	assigner Assigner
}

func NewService(topicAssigner Assigner, topicDeleter Deleter) *Service {
	return &Service{
		assigner: topicAssigner,
		deleter:  topicDeleter,
	}
}

type Deleter interface {
	DeleteTopic(ctx context.Context, topicID topic.ID, market market.Name) error
}

type Assigner interface {
	AssignToBook(ctx context.Context, topicID topic.ID, bookID book.ID, market market.Name) error
}
