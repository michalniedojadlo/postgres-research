package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// there could be an interface type e.g. DatabaseOperator
//type DatabaseOperator interface {
//	GetConn(ctx context.Context) (*pgxpool.Conn, error)
//	Ping(ctx context.Context) error
//	Disconnect()
//}

type Client struct {
	connectionPool *pgxpool.Pool
}

// Connect attempts to create a connection to the database designated in connStr
func Connect(connStr string) (*Client, error) {
	connectionPool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres connection pool failed: %w", err)
	}

	err = connectionPool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ensuring connection is stable failed: %w", err)
	}

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

// Ping could be renamed to Check to satisfy Mirko's health check interface
func (client *Client) Ping(ctx context.Context) error {
	return client.connectionPool.Ping(ctx)
}

func (client *Client) Disconnect() {
	client.connectionPool.Close()
}
