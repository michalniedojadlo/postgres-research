package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Client struct {
	connectionPool *pgxpool.Pool
}

func Connect(connStr string) (*Client, error) {
	connectionPool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to postgres connection connectionPool")
	}

	connectionPool.Acquire(context.Background())

	client := &Client{connectionPool: connectionPool}
	return client, nil
}

func (client *Client) GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := client.connectionPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (client *Client) Ping() error { return nil }

func (client *Client) Disconnect() error {
	return nil
}
