package main

import (
	"fmt"
	"log"

	"github.com/brainly/postgres-research/internal/infra/postgres"

	boardInfra "github.com/brainly/postgres-research/pkg/board/infra/postgres"
	bookInfra "github.com/brainly/postgres-research/pkg/book/infra/postgres"
	topicInfra "github.com/brainly/postgres-research/pkg/topic/infra/postgres"

	boardInternal "github.com/brainly/postgres-research/internal/infra/postgres/board"
	topicInternal "github.com/brainly/postgres-research/internal/infra/postgres/topic"

	"github.com/brainly/postgres-research/pkg/board"
	"github.com/brainly/postgres-research/pkg/book"
	"github.com/brainly/postgres-research/pkg/topic"
)

func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"
	client, err := postgres.Connect(connStr)
	if err != nil {
		log.Fatal(err)
	}

	boardAssigner := boardInfra.NewBoardAssigner(client)
	boardCreator := boardInfra.NewBoardCreator(client)
	boardService := board.NewService(boardAssigner, boardCreator)

	topicAssigner := topicInfra.NewTopicAssigner(client)
	topicDeleter := topicInfra.NewTopicDeleter(client)
	topicService := topic.NewService(topicAssigner, topicDeleter)

	bookCreator := bookInfra.NewBookCreator(client, boardInternal.AssignBoardsToBook, topicInternal.AssignTopicToBook)
	bookDeleter := bookInfra.NewBookDeleter(client, boardInternal.ClearBoardToBookAssignments, topicInternal.ClearTopicToBookAssignments)
	bookGetter := bookInfra.NewBookGetter(client)
	bookService := book.NewService(bookCreator, bookDeleter, bookGetter)

	fmt.Println(boardService, topicService, bookService)
}
