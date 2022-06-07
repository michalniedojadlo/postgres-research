package booktotopic

import "fmt"

const (
	TableName = `book_to_topic`

	ColumnBookID  = `book_id`
	ColumnMarket  = `market`
	ColumnTopicID = `topic_id`
)

func WithTableName(columnName string) string {
	return fmt.Sprintf(`%s.%s`, TableName, columnName)
}
