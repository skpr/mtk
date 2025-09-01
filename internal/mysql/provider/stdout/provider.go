package stdout

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/skpr/mtk/internal/mysql/provider"
	providerutils "github.com/skpr/mtk/internal/mysql/provider/utils"
)

// Client used for dumping a database and/or table.
type Client struct {
	provider.Interface
	Conn   *sql.Conn
	Logger *log.Logger
}

// NewClient for dumping a full or single table from a database.
func NewClient(conn *sql.Conn, logger *log.Logger) *Client {
	return &Client{
		Conn:   conn,
		Logger: logger,
	}
}

// WriteTableData will write the data from a table to the provided writer.
func (d *Client) WriteTableData(ctx context.Context, w io.Writer, table string, params provider.DumpParams) error {
	cols, err := providerutils.QueryColumnsForTable(ctx, d.Conn, table, params)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(cols, ", "), table)

	if where, ok := params.WhereMap[strings.ToLower(table)]; ok {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	rows, err := d.Conn.QueryContext(ctx, query)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	defer rows.Close()

	values := make([]*sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	var (
		counter  = 0
		firstRun = true
	)

	for rows.Next() {
		// We have already done a loop and need to close the previous insert statement.
		if counter >= params.ExtendedInsertRows {
			fmt.Fprintln(w, ";")
			counter = 0
		} else {
			if !firstRun {
				fmt.Fprint(w, ",")
			}
		}

		if counter == 0 {
			fmt.Fprintf(w, "INSERT INTO `%s` VALUES ", table)
		}

		counter++

		firstRun = false

		if err = rows.Scan(scanArgs...); err != nil {
			return err
		}

		var vals []string

		for _, col := range values {
			val := "NULL"

			if col != nil {
				val, err = getValue(string(*col))
				if err != nil {
					return err
				}
			}

			vals = append(vals, val)
		}

		fmt.Fprintf(w, "(%s)", strings.Join(vals, ","))
	}

	fmt.Fprintln(w, ";")

	return nil
}
