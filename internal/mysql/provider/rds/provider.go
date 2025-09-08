package rds

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

	Region string // Region configuration
	URI    string // S3 URI configuration
}

// NewClient for dumping a full or single table from a database.
func NewClient(conn *sql.Conn, logger *log.Logger, region, uri string) *Client {
	return &Client{
		Conn:   conn,
		Logger: logger,
		Region: region,
		URI:    uri,
	}
}

// WriteTableData will export the data from a table to S3 and write the LOAD DATA query to the provided writer.
func (d *Client) WriteTableData(ctx context.Context, w io.Writer, table string, params provider.DumpParams) error {
	// Push table data to s3.
	err := d.exportTableData(ctx, table, params)
	if err != nil {
		return err
	}

	// Write the import query to the writer.
	err = d.writeLoadQueryForTable(w, table)
	if err != nil {
		return err
	}

	return nil
}

// Export the data from a table to S3.
func (d *Client) exportTableData(ctx context.Context, table string, params provider.DumpParams) error {
	d.Logger.Printf("Exporting data to S3 for table: %s", table)

	cols, err := providerutils.QueryColumnsForTable(ctx, d.Conn, table, params)
	if err != nil {
		return err
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

	_, err = d.Conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error exporting data to S3 for table %s: %w", table, err)
	}

	return nil
}

// Write a LOAD DATA FROM S3 query to the provided writer.
func (d *Client) writeLoadQueryForTable(w io.Writer, table string) error {
	if table == "" {
		return fmt.Errorf("error: no table specified")
	}

	if d.Region == "" || len(strings.Split(d.Region, "-")) != 3 {
		return fmt.Errorf("error: region is not configured correctly")
	}

	path := strings.TrimPrefix(d.URI, "s3://")
	query := fmt.Sprintf("LOAD DATA FROM S3 MANIFEST 'S3-%s://%s/%s.csv.manifest' INTO TABLE `%s` CHARACTER SET utf8mb4", d.Region, path, table, table)
	query = fmt.Sprintf("%s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n'", query)

	_, err := fmt.Fprintln(w, query)
	return err
}
