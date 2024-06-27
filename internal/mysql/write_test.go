package mysql

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/mtk/internal/mysql/mock"
	"github.com/skpr/mtk/internal/mysql/provider"
)

func TestMySQLDumpTableHeader(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "stdout", "", "")
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `table`").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1234))
	buffer := bytes.NewBuffer(make([]byte, 0))
	count, err := dumper.WriteTableHeader(buffer, "table", provider.DumpParams{})
	assert.Equal(t, uint64(1234), count)
	assert.Nil(t, err)
	assert.Contains(t, buffer.String(), "Data for table `table`")
	assert.Contains(t, buffer.String(), "1234 rows")
}

func TestMySQLDumpTableHeaderHandlingError(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "stdout", "", "")
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `table`").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(nil))
	buffer := bytes.NewBuffer(make([]byte, 0))
	count, err := dumper.WriteTableHeader(buffer, "table", provider.DumpParams{})
	assert.Equal(t, uint64(0), count)
	assert.NotNil(t, err)
}

func TestMySQLDumpTableLockWrite(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	dumper := NewClient(nil, log.New(os.Stdout, "", 0), "stdout", "", "")
	dumper.WriteTableLockWrite(buffer, "table")
	assert.Contains(t, buffer.String(), "LOCK TABLES `table` WRITE;")
}

func TestMySQLDumpUnlockTables(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	dumper := NewClient(nil, log.New(os.Stdout, "", 0), "stdout", "", "")
	dumper.WriteUnlockTables(buffer)
	assert.Contains(t, buffer.String(), "UNLOCK TABLES;")
}

func TestMySQLDumpTableData(t *testing.T) {
	db, mock := mock.GetDB(t)
	buffer := bytes.NewBuffer(make([]byte, 0))
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "stdout", "", "")

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

	assert.Nil(t, dumper.WriteTableData(buffer, "table", provider.DumpParams{
		ExtendedInsertRows: 2}))

	assert.Equal(t, strings.Count(buffer.String(), "INSERT INTO `table` VALUES"), 3)
	assert.Equal(t, buffer.String(), "INSERT INTO `table` VALUES (1,'Go'),(2,'Java');\nINSERT INTO `table` VALUES (3,'C'),(4,'C++');\nINSERT INTO `table` VALUES (5,'Rust'),(6,'Closure');\n")
}

func TestMySQLDumpTableDataHandlingErrorFromSelectAllDataFor(t *testing.T) {
	db, mock := mock.GetDB(t)
	buffer := bytes.NewBuffer(make([]byte, 0))
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "stdout", "", "")
	error := errors.New("fail")
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnError(error)
	assert.Equal(t, error, dumper.WriteTableData(buffer, "table", provider.DumpParams{}))
}
