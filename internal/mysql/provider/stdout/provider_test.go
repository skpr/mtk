package stdout

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/mtk/internal/mysql/mock"
	"github.com/skpr/mtk/internal/mysql/provider"
)

func TestMySQLGetSelectQueryFor(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnRows(
		sqlmock.NewRows([]string{"c1", "c2"}).AddRow("a", "b"))
	query, err := dumper.GetSelectQueryForTable("table", provider.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"c2": "NOW()"}},
		WhereMap:  map[string]string{"table": "c1 > 0"},
	})
	assert.Nil(t, err)
	assert.Equal(t, "SELECT `c1`, NOW() AS `c2` FROM `table` WHERE c1 > 0", query)
}

func TestMySQLGetSelectQueryForHandlingError(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	e := errors.New("broken")
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnError(e)
	query, err := dumper.GetSelectQueryForTable("table", provider.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"c2": "NOW()"}},
		WhereMap:  map[string]string{"table": "c1 > 0"},
	})
	assert.Equal(t, e, err)
	assert.Equal(t, "", query)
}
