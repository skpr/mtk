package rds

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/skpr/mtk/internal/mysql/provider"
	providerutils "github.com/skpr/mtk/internal/mysql/provider/utils"
)

// Client used for dumping a database and/or table.
type Client struct {
	provider.Interface
	DB     *sql.DB
	Logger *log.Logger

	Region string // Region configuration
	URI    string // S3 URI configuration
}

// NewClient for dumping a full or single table from a database.
func NewClient(db *sql.DB, logger *log.Logger, region, uri string) *Client {
	return &Client{
		DB:     db,
		Logger: logger,
		Region: region,
		URI:    uri,
	}
}

// GetSelectQueryForTable will return a complete SELECT query to export data from a table.
func (d *Client) GetSelectQueryForTable(table string, params provider.DumpParams) (string, error) {
	cols, err := providerutils.QueryColumnsForTable(d.DB, table, params)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("SELECT %s", strings.Join(cols, ", "))
	query = fmt.Sprintf("%s FROM `%s`", query, table)

	if where, ok := params.WhereMap[strings.ToLower(table)]; ok {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	query = fmt.Sprintf("%s INTO OUTFILE S3 '%s/%s.csv'", query, d.URI, table)
	query = fmt.Sprintf("%s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n'", query)
	query = fmt.Sprintf("%s MANIFEST ON", query)
	query = fmt.Sprintf("%s OVERWRITE ON", query)

	importQuery, err := d.GetLoadQueryForTable(table)
	if err != nil {
		return "", err
	}

	fmt.Println(importQuery)
	return query, nil
}

// GetLoadQueryForTable will return a complete SELECT query to fetch data from a table.
func (d *Client) GetLoadQueryForTable(table string) (string, error) {
	if table == "" {
		return "", fmt.Errorf("error: no table specified")
	}
	if d.Region == "" || len(strings.Split(d.Region, "-")) != 3 {
		return "", fmt.Errorf("error: region is not configured correctly")
	}
	path := strings.TrimPrefix(d.URI, "s3://")
	query := fmt.Sprintf("LOAD DATA FROM S3 MANIFEST 'S3-%s://%s/%s.csv.manifest' INTO TABLE `%s`", d.Region, path, table, table)
	query = fmt.Sprintf("%s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n'", query)

	return query, nil
}
