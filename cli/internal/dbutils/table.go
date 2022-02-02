package dbutils

import (
	"database/sql"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/skpr/mtk/internal/sliceutils"
)

// ListTables based on a set of globs.
func ListTables(db *sql.DB, globs []string) ([]string, error) {
	var globbed []string

	tables, err := queryTables(db)
	if err != nil {
		return globbed, errors.Wrap(err, "failed to query for tables")
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

// Helper function to get a list of tables.
func queryTables(db *sql.DB) ([]string, error) {
	var tables []string

	rows, err := db.Query("SHOW FULL TABLES")
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
