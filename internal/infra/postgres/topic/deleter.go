package topic

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/brainly/postgres-research/internal/core/market"
	"github.com/brainly/postgres-research/internal/core/topic"
	"github.com/brainly/postgres-research/internal/infra/postgres"

	topicSchema "github.com/brainly/postgres-research/internal/infra/postgres/schema/topic"
)

func DeleteTopic(ctx context.Context, topicID topic.ID, market market.Name, queryBuilder squirrel.StatementBuilderType, executioner postgres.QueryExecutioner) error {
	query, args, err := queryBuilder.
		Delete(topicSchema.TableName).
		Where(squirrel.Eq{
			topicSchema.ColumnID:     topicID,
			topicSchema.ColumnMarket: market,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("creating query failed: %w", err)
	}

	_, err = executioner.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("executing query failed: %w", err)
	}

	return nil
}
