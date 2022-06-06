package board

import (
	"github.com/brainly/postgres-research/internal/core/board"
	"github.com/brainly/postgres-research/internal/core/book"
	"github.com/brainly/postgres-research/internal/core/market"
)

type CommandAssignBookToBoard struct {
	BookID   book.ID
	BoardIDs []board.ID
	Market   market.Name
}

type CommandCreateBoard struct {
	Market market.Name
	Name   string
}
