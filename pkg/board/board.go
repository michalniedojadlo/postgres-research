package board

import (
	"context"
	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
)

type Board struct {
	ID     board.ID
	Name   string
	Market market.Name
}

type Service struct {
	assigner Assigner
	creator  Creator
}

func NewService(boardAssigner Assigner, boardCreator Creator) *Service {
	return &Service{
		assigner: boardAssigner,
		creator:  boardCreator,
	}
}

type Assigner interface {
	AssignToBook(ctx context.Context, boardIDs board.IDs, bookID book.ID, market market.Name) error
}

type Creator interface {
	Create(ctx context.Context, id board.ID, market market.Name, name string) error
}
