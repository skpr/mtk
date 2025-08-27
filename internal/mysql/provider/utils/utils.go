package utils

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/skpr/mtk/internal/mysql/provider"
)

// QueryColumnsForTable for a given table.
func QueryColumnsForTable(ctx context.Context, conn *sql.Conn, table string, params provider.DumpParams) ([]string, error) {
	var rows *sql.Rows

	rows, err := conn.QueryContext(ctx, fmt.Sprintf("SELECT * FROM `%s` LIMIT 1", table))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for k, column := range columns {
		replacement, ok := params.SelectMap[strings.ToLower(table)][strings.ToLower(column)]
		if ok {
			columns[k] = fmt.Sprintf("%s AS `%s`", replacement, column)
		} else {
			columns[k] = fmt.Sprintf("`%s`", column)
		}
	}

	return columns, nil
}
