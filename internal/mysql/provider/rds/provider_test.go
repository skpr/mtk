package rds

import (
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/mtk/internal/mysql/mock"
	"github.com/skpr/mtk/internal/mysql/provider"
)

func TestMySQLGetExportSelectQueryFor(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "ap-southheast-2", "s3://path/to/bucket")
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnRows(
		sqlmock.NewRows([]string{"c1", "c2"}).AddRow("a", "b"))
	query, err := dumper.GetSelectQueryForTable("table", provider.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"c2": "NOW()"}},
		WhereMap:  map[string]string{"table": "c1 > 0"},
	})
	assert.Nil(t, err)
	assert.Equal(t, "SELECT `c1`, NOW() AS `c2` FROM `table` WHERE c1 > 0 INTO OUTFILE S3 's3://path/to/bucket/table.csv' FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n' MANIFEST ON OVERWRITE ON", query)
}

func TestMySQLGetLoadQueryFor(t *testing.T) {
	db, _ := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0), "ap-southeast-4", "s3://path/to/bucket")
	query, err := dumper.GetLoadQueryForTable("table_name")
	assert.Nil(t, err)
	assert.Equal(t, "LOAD DATA FROM S3 MANIFEST 'S3-ap-southeast-4://path/to/bucket/table_name.csv.manifest' INTO TABLE `table_name` FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n'", query)

}
