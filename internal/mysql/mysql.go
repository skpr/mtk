package mysql

import (
	"database/sql"
	"fmt"
	"io"
	"log"

	"github.com/go-sql-driver/mysql"
)

const (
	// OperationIgnore is used to skip a table when dumping.
	OperationIgnore = "ignore"
	// OperationNoData is used when you want to dump a table structure without the data.
	OperationNoData = "nodata"
)

type Connection struct {
	Hostname string
	Username string
	Password string
	Protocol string
	Port     int32
	MaxConn  int
}

func (o Connection) Open(database string) (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 o.Username,
		Passwd:               o.Password,
		Net:                  o.Protocol,
		Addr:                 fmt.Sprintf("%s:%d", o.Hostname, o.Port),
		DBName:               database,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(o.MaxConn)

	return db, nil
}

// Client used for dumping a database and/or table.
type Client struct {
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

// DumpParams is used to pass parameters to the Dump function.
type DumpParams struct {
	SelectMap          map[string]map[string]string
	WhereMap           map[string]string
	FilterMap          map[string]string
	UseTableLock       bool
	ExtendedInsertRows int
}

// DumpTables will write all table data to a single writer.
func (d *Client) DumpTables(w io.Writer, params DumpParams) error {
	if err := d.WriteHeader(w); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if err := d.writeTables(w, params); err != nil {
		return fmt.Errorf("failed to write tables: %w", err)
	}

	if err := d.WriteFooter(w); err != nil {
		return fmt.Errorf("failed to write footer: %w", err)
	}

	if err := d.WriteDumpCompleted(w); err != nil {
		return fmt.Errorf("failed to write completed datetime: %w", err)
	}

	return nil
}

// DumpTable is convenient if you wish to coordinate a dump eg. Single file per table.
func (d *Client) DumpTable(w io.Writer, table string, params DumpParams) error {
	if err := d.WriteHeader(w); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if err := d.writeTable(w, table, params); err != nil {
		return fmt.Errorf("failed to write tables: %w", err)
	}

	if err := d.WriteFooter(w); err != nil {
		return fmt.Errorf("failed to write footer: %w", err)
	}

	if err := d.WriteDumpCompleted(w); err != nil {
		return fmt.Errorf("failed to write completed datetime: %w", err)
	}

	return nil
}
