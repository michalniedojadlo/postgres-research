package board

import (
	"context"
)

func (svc *Service) AssignToBook(ctx context.Context, cmd CommandAssignBookToBoard) error {
	err := svc.assigner.AssignToBook(ctx, cmd.BoardIDs, cmd.BookID, cmd.Market)
	if err != nil {
		return err
	}

	return nil
}
