package board

import (
	"context"
	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/google/uuid"
)

func (svc *Service) CreateBoard(ctx context.Context, cmd CommandCreateBoard) error {
	newBoardID := uuid.New()
	err := svc.creator.Create(ctx, board.ID(newBoardID), cmd.Market, cmd.Name)
	if err != nil {
		return err
	}

	return nil
}
