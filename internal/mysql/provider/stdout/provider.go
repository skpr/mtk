package stdout

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
}

// NewClient for dumping a full or single table from a database.
func NewClient(db *sql.DB, logger *log.Logger) *Client {
	return &Client{
		DB:     db,
		Logger: logger,
	}
}

// GetSelectQueryForTable will return a complete SELECT query to fetch data from a table.
func (d *Client) GetSelectQueryForTable(table string, params provider.DumpParams) (string, error) {
	cols, err := providerutils.QueryColumnsForTable(d.DB, table, params)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(cols, ", "), table)

	if where, ok := params.WhereMap[strings.ToLower(table)]; ok {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	return query, nil
}
