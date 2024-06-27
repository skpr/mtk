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

	DataPath string // S3 URI configuration
	Region   string // Region configuration
}

// NewClient for dumping a full or single table from a database.
func NewClient(db *sql.DB, logger *log.Logger) *Client {
	return &Client{
		DB:     db,
		Logger: logger,
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

	query = fmt.Sprintf("%s INTO OUTFILE S3 '%s/%s.csv'", query, d.DataPath, table)
	query = fmt.Sprintf("%s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n'", query)
	query = fmt.Sprintf("%s MANIFEST ON", query)
	query = fmt.Sprintf("%s OVERWRITE ON", query)

	importQuery, err := d.GetLoadQueryForTable(table, d.DataPath, d.Region)
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
	query := fmt.Sprintf("LOAD DATA FROM S3 MANIFEST 'S3-%s://%s/%s.csv.manifest' INTO TABLE `%s`", region, path, table, table)
	query = fmt.Sprintf("%s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n'", query)

	return query, nil
}
