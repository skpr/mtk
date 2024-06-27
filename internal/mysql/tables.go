package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gobwas/glob"

	"github.com/skpr/mtk/internal/mysql/provider"
	"github.com/skpr/mtk/internal/mysql/provider/rds"
	"github.com/skpr/mtk/internal/mysql/provider/stdout"
	"github.com/skpr/mtk/internal/sliceutils"
)

// ListTablesByGlob will return a list of tables based on a list of globs.
func (d *Client) ListTablesByGlob(globs []string) ([]string, error) {
	var globbed []string

	tables, err := d.QueryTables()
	if err != nil {
		return globbed, fmt.Errorf("failed to query for tables: %w", err)
	}

	for _, query := range globs {
		g := glob.MustCompile(query)

		for _, table := range tables {
			if g.Match(table) {
				globbed = sliceutils.AppendIfMissing(globbed, table)
			}
		}
	}

	return globbed, nil
}

// QueryTables will return a list of tables.
func (d *Client) QueryTables() ([]string, error) {
	// Use the cached tables if we have them.
	if len(d.cachedTables) > 0 {
		return d.cachedTables, nil
	}

	tables := make([]string, 0)

	rows, err := d.DB.Query("SHOW FULL TABLES")
	if err != nil {
		return tables, err
	}

	defer rows.Close()

	for rows.Next() {
		var tableName, tableType string

		err := rows.Scan(&tableName, &tableType)
		if err != nil {
			return tables, err
		}

		if tableType == "BASE TABLE" {
			tables = append(tables, tableName)
		}
	}

	// Set the cached tables for future executions.
	d.cachedTables = tables

	return tables, nil
}

func (d *Client) getProviderClient() (provider.Interface, error) {
	switch d.Provider {
	case "rds":
		client := rds.NewClient(d.DB, d.Logger, d.Region, d.URI)
		return client, nil
	case "stdout":
		return stdout.NewClient(d.DB, d.Logger), nil
	default:
		return nil, errors.New("invalid provider")
	}
}

// Helper function to get all data for a table.
func (d *Client) selectAllDataForTable(table string, params provider.DumpParams) (*sql.Rows, []string, error) {

	client, err := d.getProviderClient()
	if err != nil {
		return nil, nil, err
	}

	query, err := client.GetSelectQueryForTable(table, params)
	if err != nil {
		return nil, nil, err
	}

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	return rows, columns, nil
}

// GetRowCountForTable will return the number of rows using a SELECT statement.
func (d *Client) GetRowCountForTable(table string, params provider.DumpParams) (uint64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM `%s`", table)

	if where, ok := params.WhereMap[strings.ToLower(table)]; ok {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	row := d.DB.QueryRow(query)

	var count uint64

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// LockTableReading explicitly acquires table locks for the current client session.
func (d *Client) LockTableReading(table string) (sql.Result, error) {
	return d.DB.Exec(fmt.Sprintf("LOCK TABLES `%s` READ", table))
}

// UnlockTables explicitly releases any table locks held by the current session.
func (d *Client) UnlockTables() (sql.Result, error) {
	return d.DB.Exec("UNLOCK TABLES")
}

// FlushTable will force a tables to be closed.
func (d *Client) FlushTable(table string) (sql.Result, error) {
	return d.DB.Exec(fmt.Sprintf("FLUSH TABLES `%s`", table))
}
