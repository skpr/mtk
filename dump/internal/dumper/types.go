package dumper

// TableType used to identify what type the table is.
type TableType string

const (
	// TypeTableBase are the underlying tables that actually store the metadata for a specific database.
	TypeTableBase TableType = "BASE TABLE"
	// TypeTableView is a virtual table based on the result-set of an SQL statement.
	TypeTableView TableType = "VIEW"
)

// Table is a data set consisting of columns and rows.
type Table struct {
	Name string
	Type TableType
}
