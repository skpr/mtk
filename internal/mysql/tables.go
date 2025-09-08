package mysql

import (
	"context"
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
func (d *Client) ListTablesByGlob(ctx context.Context, globs []string) ([]string, error) {
	var globbed []string

	tables, err := d.QueryTables(ctx)
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
func (d *Client) QueryTables(ctx context.Context) ([]string, error) {
	// Use the cached tables if we have them.
	if len(d.cachedTables) > 0 {
		return d.cachedTables, nil
	}

	tables := make([]string, 0)

	rows, err := d.Conn.QueryContext(ctx, "SHOW FULL TABLES")
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
		client := rds.NewClient(d.Conn, d.Logger, d.Region, d.URI)
		return client, nil
	case "stdout":
		return stdout.NewClient(d.Conn, d.Logger), nil
	default:
		return nil, errors.New("invalid provider")
	}
}

// GetRowCountForTable will return the number of rows using a SELECT statement.
func (d *Client) GetRowCountForTable(ctx context.Context, table string, params provider.DumpParams) (uint64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM `%s`", table)

	if where, ok := params.WhereMap[strings.ToLower(table)]; ok {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	row := d.Conn.QueryRowContext(ctx, query)

	var count uint64

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// LockTableReading explicitly acquires table locks for the current client session.
func (d *Client) LockTableReading(ctx context.Context, table string) (sql.Result, error) {
	return d.Conn.ExecContext(ctx, fmt.Sprintf("LOCK TABLES `%s` READ", table))
}

// UnlockTables explicitly releases any table locks held by the current session.
func (d *Client) UnlockTables(ctx context.Context) (sql.Result, error) {
	return d.Conn.ExecContext(ctx, "UNLOCK TABLES")
}

// FlushTable will force a tables to be closed.
func (d *Client) FlushTable(ctx context.Context, table string) (sql.Result, error) {
	return d.Conn.ExecContext(ctx, fmt.Sprintf("FLUSH TABLES `%s`", table))
}
