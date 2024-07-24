package utils

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/mtk/internal/mysql/mock"
	"github.com/skpr/mtk/internal/mysql/provider"
)

func TestMySQLGetColumnsForSelect(t *testing.T) {
	db, mock := mock.GetDB(t)
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnRows(
		sqlmock.NewRows([]string{"col1", "col2", "col3"}).AddRow("a", "b", "c"))
	columns, err := QueryColumnsForTable(db, "table", provider.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"col2": "NOW()"}},
	})
	assert.Nil(t, err)
	assert.Equal(t, []string{"`col1`", "NOW() AS `col2`", "`col3`"}, columns)
}

func TestMySQLGetColumnsForSelectHandlingErrorWhenQuerying(t *testing.T) {
	db, mock := mock.GetDB(t)
	error := errors.New("broken")
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnError(error)
	columns, err := QueryColumnsForTable(db, "table", provider.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"col2": "NOW()"}},
	})
	assert.Equal(t, err, error)
	assert.Empty(t, columns)
}
