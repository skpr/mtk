package rds

import (
	"bytes"
	"context"
	"io"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/skpr/mtk/internal/mysql/provider"
)

func TestWriteTableData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("db.Conn error: %v", err)
	}

	// 1) QueryColumnsForTable runs: expect the probe query and return columns.
	mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` LIMIT 1")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"})) // columns only is fine

	// 2) exportTableData runs: now that code uses ExecContext, expect an Exec.
	exportSQL := "SELECT `id`, `name` FROM `users` INTO OUTFILE S3 's3://my-bucket/prefix/users.csv' " +
		"FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n' MANIFEST ON OVERWRITE ON"
	mock.
		ExpectExec(regexp.QuoteMeta(exportSQL)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // no rows affected is fine

	// Client + buffer for the LOAD DATA query
	buf := &bytes.Buffer{}
	c := &Client{
		Conn:   conn,
		Logger: log.New(io.Discard, "", 0),
		Region: "ap-southeast-2",
		URI:    "s3://my-bucket/prefix",
	}

	// Run the full flow
	if err := c.WriteTableData(context.Background(), buf, "users", provider.DumpParams{
		WhereMap:  map[string]string{},
		SelectMap: map[string]map[string]string{},
	}); err != nil {
		t.Fatalf("WriteTableData error: %v", err)
	}

	// Assert the LOAD DATA query is exactly written (including newline)
	got := buf.String()
	want := "LOAD DATA FROM S3 MANIFEST 'S3-ap-southeast-2://my-bucket/prefix/users.csv.manifest' " +
		"INTO TABLE `users` CHARACTER SET utf8mb4 FIELDS TERMINATED BY ',' ENCLOSED BY '\"' " +
		"LINES TERMINATED BY '\\n'\n"
	if got != want {
		t.Fatalf("unexpected load query:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}

	// Ensure all expectations met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
