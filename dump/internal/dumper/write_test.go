package dumper

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/mtk/dump/internal/dumper/mock"
)

func TestMySQLDumpTableHeader(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `table`").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1234))
	buffer := bytes.NewBuffer(make([]byte, 0))
	count, err := dumper.WriteTableHeader(buffer, "table")
	assert.Equal(t, uint64(1234), count)
	assert.Nil(t, err)
	assert.Contains(t, buffer.String(), "Data for table `table`")
	assert.Contains(t, buffer.String(), "1234 rows")
}

func TestMySQLDumpTableHeaderHandlingError(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `table`").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(nil))
	buffer := bytes.NewBuffer(make([]byte, 0))
	count, err := dumper.WriteTableHeader(buffer, "table")
	assert.Equal(t, uint64(0), count)
	assert.NotNil(t, err)
}

func TestMySQLDumpTableLockWrite(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	dumper := NewClient(nil, log.New(os.Stdout, "", 0))
	dumper.WriteTableLockWrite(buffer, "table")
	assert.Contains(t, buffer.String(), "LOCK TABLES `table` WRITE;")
}

func TestMySQLDumpUnlockTables(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	dumper := NewClient(nil, log.New(os.Stdout, "", 0))
	dumper.WriteUnlockTables(buffer)
	assert.Contains(t, buffer.String(), "UNLOCK TABLES;")
}

func TestMySQLDumpTableData(t *testing.T) {
	db, mock := mock.GetDB(t)
	buffer := bytes.NewBuffer(make([]byte, 0))
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	dumper.ExtendedInsertRows = 2

	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnRows(
		sqlmock.NewRows([]string{"id", "language"}).
			AddRow(1, "Go"))

	mock.ExpectQuery("SELECT `id`, `language` FROM `table`").WillReturnRows(
		sqlmock.NewRows([]string{"id", "language"}).
			AddRow(1, "Go").
			AddRow(2, "Java").
			AddRow(3, "C").
			AddRow(4, "C++").
			AddRow(5, "Rust").
			AddRow(6, "Closure"))

	assert.Nil(t, dumper.WriteTableData(buffer, "table"))

	assert.Equal(t, strings.Count(buffer.String(), "INSERT INTO `table` VALUES"), 3)
	assert.Contains(t, buffer.String(), `'Go'`)
	assert.Contains(t, buffer.String(), `'Java'`)
	assert.Contains(t, buffer.String(), `'C'`)
	assert.Contains(t, buffer.String(), `'C++'`)
	assert.Contains(t, buffer.String(), `'Rust'`)
	assert.Contains(t, buffer.String(), `'Closure'`)
}

func TestMySQLDumpTableDataHandlingErrorFromSelectAllDataFor(t *testing.T) {
	db, mock := mock.GetDB(t)
	buffer := bytes.NewBuffer(make([]byte, 0))
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	error := errors.New("fail")
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnError(error)
	assert.Equal(t, error, dumper.WriteTableData(buffer, "table"))
}

func TestWriteCreateView(t *testing.T) {
	db, mock := mock.GetDB(t)
	buffer := bytes.NewBuffer(make([]byte, 0))
	dumper := NewClient(db, log.New(os.Stdout, "", 0))

	mock.ExpectQuery("SHOW CREATE VIEW `table`").WillReturnRows(
		sqlmock.NewRows([]string{"View", "Create View", "character_set_client", "collation_connection"}).
			AddRow(1, "CREATE VIEW test.v AS SELECT * FROM t;", "utf8mb4", "utf8mb4_0900_ai_ci"))

	err := dumper.WriteCreateView(buffer, "table")
	assert.NoError(t, err)

	want, err := os.ReadFile("testdata/view.sql")
	assert.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("\n%s\n", want), buffer.String())
}
