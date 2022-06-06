package booktoboard

const (
	TableName = `book_to_board`

	ColumnBookID  = `book_id`
	ColumnMarket  = `market`
	ColumnBoardID = `board_id`
)

func WithTableName(columnName string) string {
	return TableName + `.` + columnName
}
