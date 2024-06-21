package rds

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/skpr/mtk/internal/mysql/providers"
)

// Client used for dumping a database and/or table.
type Client struct {
	providers.MTKProvider
	DB     *sql.DB
	Logger *log.Logger
}

// NewClient for dumping a full or single table from a database.
func NewClient(db *sql.DB, logger *log.Logger) *Client {
	return &Client{
		DB:     db,
		Logger: logger,
	}
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

// GetSelectQueryForTable will return a complete SELECT query to export data from a table.
func (d *Client) GetSelectQueryForTable(table string, params providers.DumpParams) (string, error) {
	cols, err := d.QueryColumnsForTable(table, params)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("SELECT %s", strings.Join(cols, ", "))
	query = fmt.Sprintf("%s FROM `%s`", query, table)

	if where, ok := params.WhereMap[strings.ToLower(table)]; ok {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	query = fmt.Sprintf("%s INTO OUTFILE S3 '%s/%s.csv'", query, params.DataPath, table)
	query = fmt.Sprintf("%s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n'", query)
	query = fmt.Sprintf("%s MANIFEST ON", query)
	query = fmt.Sprintf("%s OVERWRITE ON", query)

	importQuery, err := d.GetLoadQueryForTable(table, params.DataPath, params.Region)
	if err != nil {
		return "", err
	}

	fmt.Println(importQuery)
	return query, nil
}

// GetLoadQueryForTable will return a complete SELECT query to fetch data from a table.
func (d *Client) GetLoadQueryForTable(table, path, region string) (string, error) {
	if table == "" {
		return "", fmt.Errorf("error: no table specified")
	}
	if region == "" || len(strings.Split(region, "-")) != 3 {
		return "", fmt.Errorf("error: region is not configured correctly")
	}
	path = strings.TrimPrefix(path, "s3://")
	query := fmt.Sprintf("LOAD DATA FROM S3 FILE 'S3-%s://%s/%s.csv.manifest' INTO TABLE `%s`", region, path, table, table)
	query = fmt.Sprintf("%s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n'", query)

	return query, nil
}
