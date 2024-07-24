package mysql

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/mtk/internal/mysql/mock"
	"github.com/skpr/mtk/internal/mysql/provider"
)

func TestMySQLFlushTable(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	mock.ExpectExec("FLUSH TABLES `table`").WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := dumper.FlushTable("table")
	assert.Nil(t, err)
}

func TestMySQLUnlockTables(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	mock.ExpectExec("UNLOCK TABLES").WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := dumper.UnlockTables()
	assert.Nil(t, err)
}

func TestMySQLQueryTables(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	mock.ExpectQuery("SHOW FULL TABLES").WillReturnRows(
		sqlmock.NewRows([]string{"Tables_in_database", "Table_type"}).
			AddRow("table1", "BASE TABLE").
			AddRow("table2", "BASE TABLE"),
	)
	tables, err := dumper.QueryTables()
	assert.Equal(t, []string{"table1", "table2"}, tables)
	assert.Nil(t, err)
}

func TestMySQLLockTableRead(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	mock.ExpectExec("LOCK TABLES `table` READ").WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := dumper.LockTableReading("table")
	assert.Nil(t, err)
}

func TestMySQLGetTablesHandlingErrorWhenListingTables(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	expectedErr := errors.New("broken")
	mock.ExpectQuery("SHOW FULL TABLES").WillReturnError(expectedErr)
	tables, err := dumper.QueryTables()
	assert.Equal(t, []string{}, tables)
	assert.Equal(t, expectedErr, err)
}

func TestMySQLGetTablesHandlingErrorWhenScanningRow(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	mock.ExpectQuery("SHOW FULL TABLES").WillReturnRows(
		sqlmock.NewRows([]string{"Tables_in_database", "Table_type"}).AddRow(1, nil))
	tables, err := dumper.QueryTables()
	assert.Equal(t, []string{}, tables)
	assert.NotNil(t, err)
}

func TestMySQLDumpCreateTable(t *testing.T) {
	var ddl = "CREATE TABLE `table` (" +
		"`id` bigint(20) NOT NULL AUTO_INCREMENT, " +
		"`name` varchar(255) NOT NULL, " +
		"PRIMARY KEY (`id`), KEY `idx_name` (`name`) " +
		") ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8"
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	mock.ExpectQuery("SHOW CREATE TABLE `table`").WillReturnRows(
		sqlmock.NewRows([]string{"Table", "Create Table"}).
			AddRow("table", ddl),
	)
	buffer := bytes.NewBuffer(make([]byte, 0))
	assert.Nil(t, dumper.WriteCreateTable(buffer, "table"))
	assert.Contains(t, buffer.String(), "DROP TABLE IF EXISTS `table`")
	assert.Contains(t, buffer.String(), ddl)
}

func TestMySQLDumpCreateTableHandlingErrorWhenScanningRows(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	mock.ExpectQuery("SHOW CREATE TABLE `table`").WillReturnRows(
		sqlmock.NewRows([]string{"Table", "Create Table"}).AddRow("table", nil))
	buffer := bytes.NewBuffer(make([]byte, 0))
	assert.NotNil(t, dumper.WriteCreateTable(buffer, "table"))
}

func TestMySQLGetRowCount(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `table` WHERE c1 > 0").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1234))
	count, err := dumper.GetRowCountForTable("table", provider.DumpParams{
		WhereMap: map[string]string{"table": "c1 > 0"},
	})
	assert.Nil(t, err)
	assert.Equal(t, uint64(1234), count)
}

func TestMySQLGetRowCountHandlingError(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "", "", "")
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `table` WHERE c1 > 0").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(nil))
	count, err := dumper.GetRowCountForTable("table", provider.DumpParams{
		WhereMap: map[string]string{"table": "c1 > 0"},
	})
	assert.NotNil(t, err)
	assert.Equal(t, uint64(0), count)
}
