package topic

import (
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
)

type DeleteTopicCMD struct {
	ID     topic.ID
	Market market.Name
}

type AssignTopicToBookCMD struct {
	TopicID topic.ID
	BookID  book.ID
	Market  market.Name
}
