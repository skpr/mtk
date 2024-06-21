package stdout

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/skpr/mtk/internal/mysql/mock"
	"github.com/skpr/mtk/internal/mysql/providers"
	"github.com/stretchr/testify/assert"
)

func TestMySQLGetSelectQueryFor(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnRows(
		sqlmock.NewRows([]string{"c1", "c2"}).AddRow("a", "b"))
	query, err := dumper.GetSelectQueryForTable("table", providers.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"c2": "NOW()"}},
		WhereMap:  map[string]string{"table": "c1 > 0"},
	})
	assert.Nil(t, err)
	assert.Equal(t, "SELECT `c1`, NOW() AS `c2` FROM `table` WHERE c1 > 0", query)
}

func TestMySQLGetSelectQueryForHandlingError(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	error := errors.New("broken")
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnError(error)
	query, err := dumper.GetSelectQueryForTable("table", providers.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"c2": "NOW()"}},
		WhereMap:  map[string]string{"table": "c1 > 0"},
	})
	assert.Equal(t, error, err)
	assert.Equal(t, "", query)
}

func TestMySQLGetColumnsForSelect(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnRows(
		sqlmock.NewRows([]string{"col1", "col2", "col3"}).AddRow("a", "b", "c"))
	columns, err := dumper.QueryColumnsForTable("table", providers.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"col2": "NOW()"}},
	})
	assert.Nil(t, err)
	assert.Equal(t, []string{"`col1`", "NOW() AS `col2`", "`col3`"}, columns)
}

func TestMySQLGetColumnsForSelectHandlingErrorWhenQuerying(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	error := errors.New("broken")
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnError(error)
	columns, err := dumper.QueryColumnsForTable("table", providers.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"col2": "NOW()"}},
	})
	assert.Equal(t, err, error)
	assert.Empty(t, columns)
}
