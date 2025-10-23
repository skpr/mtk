package stdout

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

	// Probe from utils.QueryColumnsForTable
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` LIMIT 1")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	// Main SELECT that GetSelectQueryForTable issues
	mainSelect := "SELECT `id`, `name` FROM `users`"
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob").
		AddRow(3, "Charlie").
		AddRow(4, "Dana").
		AddRow(5, "Eve")
	mock.ExpectQuery(regexp.QuoteMeta(mainSelect)).WillReturnRows(rows)

	var out bytes.Buffer
	c := &Client{
		Conn:   conn,
		Logger: log.New(io.Discard, "", 0),
	}

	params := provider.DumpParams{
		ExtendedInsertRows: 2,
		WhereMap:           map[string]string{},
		SelectMap:          map[string]map[string]string{},
	}

	if err := c.WriteTableData(context.Background(), &out, "users", params); err != nil {
		t.Fatalf("WriteTableData returned error: %v", err)
	}

	got := out.String()

	// Expect two batches of 2 and one batch of 1 (with final semicolon inline).
	want :=
		"INSERT INTO `users` VALUES (1,'Alice'),(2,'Bob');\n" +
			"INSERT INTO `users` VALUES (3,'Charlie'),(4,'Dana');\n" +
			"INSERT INTO `users` VALUES (5,'Eve');\n"

	if got != want {
		t.Fatalf("unexpected writer output:\n--- got ---\n%s--- want ---\n%s", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
