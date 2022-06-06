package postgres

import "fmt"

func NewErrAcquiringConnection(err error) error {
	return fmt.Errorf("acquiring connection from the pool failed: %w", err)
}
