package topic

import "context"

func (svc *Service) DeleteTopic(ctx context.Context, cmd DeleteTopicCMD) error {
	err := svc.deleter.DeleteTopic(ctx, cmd.ID, cmd.Market)
	if err != nil {
		return err
	}

	return nil
}
