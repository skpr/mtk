package provider

import (
	"context"
	"io"
)

// Interface implements the required functionality for a Provider.
type Interface interface {
	WriteTableData(ctx context.Context, w io.Writer, table string, params DumpParams) error
}

// DumpParams is used to pass parameters to the Dump function.
type DumpParams struct {
	SelectMap          map[string]map[string]string
	WhereMap           map[string]string
	FilterMap          map[string]string
	UseTableLock       bool
	ExtendedInsertRows int
}
