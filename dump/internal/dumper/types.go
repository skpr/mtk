package dumper

type TableType string

const (
	TypeTableBase TableType = "BASE TABLE"
	TypeTableView TableType = "VIEW"
)

type Table struct {
	Name string
	Type TableType
}
