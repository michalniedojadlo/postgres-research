package topic

import "context"

func (svc *Service) AssignToBook(ctx context.Context, cmd AssignTopicToBookCMD) error {
	err := svc.assigner.AssignToBook(ctx, cmd.TopicID, cmd.BookID, cmd.Market)
	if err != nil {
		return err
	}

	return nil
}
