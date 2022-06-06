package book

import (
	"context"
)

func (service *Service) DeleteBook(ctx context.Context, cmd *DeleteBookCMD) error {
	err := service.bookDeleter.DeleteBookByID(ctx, cmd.ID, cmd.Market)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) SoftDelete(ctx context.Context, cmd *DeleteBookCMD) error {
	err := service.bookDeleter.DeleteBookByID(ctx, cmd.ID, cmd.Market)
	if err != nil {
		return err
	}

	return nil
}
