package dumper

import (
	"database/sql"
	"fmt"
	"strings"
)

// QueryTables will return a list of tables.
func (d *Client) QueryTables() ([]string, error) {
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

	return tables, nil
}

// QueryColumnsForTable for a given table.
func (d *Client) QueryColumnsForTable(table string) ([]string, error) {
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
		replacement, ok := d.SelectMap[strings.ToLower(table)][strings.ToLower(column)]
		if ok {
			columns[k] = fmt.Sprintf("%s AS `%s`", replacement, column)
		} else {
			columns[k] = fmt.Sprintf("`%s`", column)
		}
	}

	return columns, nil
}

// GetSelectQueryForTable will return a complete SELECT query to fetch data from a table.
func (d *Client) GetSelectQueryForTable(table string) (string, error) {
	cols, err := d.QueryColumnsForTable(table)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(cols, ", "), table)

	if where, ok := d.WhereMap[strings.ToLower(table)]; ok {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	return query, nil
}

// Helper function to get all data for a table.
func (d *Client) selectAllDataForTable(table string) (*sql.Rows, []string, error) {
	query, err := d.GetSelectQueryForTable(table)
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
func (d *Client) GetRowCountForTable(table string) (uint64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM `%s`", table)

	if where, ok := d.WhereMap[strings.ToLower(table)]; ok {
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
