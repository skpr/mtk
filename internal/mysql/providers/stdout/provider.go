package stdout

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/skpr/mtk/internal/mysql/providers"
)

type DumpParams struct {
	providers.DumpParams
}

// NewClient for dumping a full or single table from a database.
func NewClient(db *sql.DB, logger *log.Logger) *Client {
	return &Client{
		DB:     db,
		Logger: logger,
	}
}

// Client used for dumping a database and/or table.
type Client struct {
	providers.MTKProvider
	DB     *sql.DB
	Logger *log.Logger
	// A field for caching a list of tables for this database.
	cachedTables []string
}

// QueryColumnsForTable for a given table.
func (d *Client) QueryColumnsForTable(table string, params providers.DumpParams) ([]string, error) {
	var rows *sql.Rows

	rows, err := d.DB.Query(fmt.Sprintf("SELECT * FROM `%s` LIMIT 1", table))
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

// GetSelectQueryForTable will return a complete SELECT query to fetch data from a table.
func (d *Client) GetSelectQueryForTable(table string, params providers.DumpParams) (string, error) {
	cols, err := d.QueryColumnsForTable(table, params)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(cols, ", "), table)

	if where, ok := params.WhereMap[strings.ToLower(table)]; ok {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	return query, nil
}
