package provider

// Interface implements the required functionality for a Provider.
type Interface interface {
	GetSelectQueryForTable(table string, params DumpParams) (string, error)
}

// DumpParams is used to pass parameters to the Dump function.
type DumpParams struct {
	SelectMap          map[string]map[string]string
	WhereMap           map[string]string
	FilterMap          map[string]string
	UseTableLock       bool
	ExtendedInsertRows int
}
