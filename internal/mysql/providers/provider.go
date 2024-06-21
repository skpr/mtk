package providers

// MTKProvider is an interface implementing the required functionality for a Provider.
type MTKProvider interface {
	QueryColumnsForTable(table string, params DumpParams) ([]string, error)
	GetSelectQueryForTable(table string, params DumpParams) (string, error)
}

// DumpParams is used to pass parameters to the Dump function.
type DumpParams struct {
	SelectMap          map[string]map[string]string
	WhereMap           map[string]string
	FilterMap          map[string]string
	UseTableLock       bool
	ExtendedInsertRows int

	Provider string // Provider configuration.
	DataPath string // S3 URI configuration for RDS provider.
	Region   string // Region configuration for RDS provider.
}
