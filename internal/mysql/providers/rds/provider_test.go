package rds

import (
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/skpr/mtk/internal/mysql/mock"
	"github.com/skpr/mtk/internal/mysql/providers"
	"github.com/stretchr/testify/assert"
)

func TestMySQLGetExportSelectQueryFor(t *testing.T) {
	db, mock := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnRows(
		sqlmock.NewRows([]string{"c1", "c2"}).AddRow("a", "b"))
	query, err := dumper.GetSelectQueryForTable("table", providers.DumpParams{
		SelectMap:  map[string]map[string]string{"table": {"c2": "NOW()"}},
		WhereMap:   map[string]string{"table": "c1 > 0"},
		DataExport: true,
		DataPath:   "s3://path/to/bucket",
		Region:     "ap-southeast-4",
	})
	assert.Nil(t, err)
	assert.Equal(t, "SELECT `c1`, NOW() AS `c2` FROM `table` WHERE c1 > 0 INTO OUTFILE S3 's3://path/to/bucket/table.csv' FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n' MANIFEST ON OVERWRITE ON", query)

}

func TestMySQLGetImportSelectQueryFor(t *testing.T) {
	db, _ := mock.GetDB(t)
	dumper := NewClient(db, log.New(os.Stdout, "", 0))
	query, err := dumper.GetLoadQueryForTable("table_name", "s3://path/to/bucket", "ap-southeast-4")
	assert.Nil(t, err)
	assert.Equal(t, "LOAD DATA FROM S3 FILE 'S3-ap-southeast-4://path/to/bucket/table_name.csv.manifest' INTO TABLE `table_name` FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\\n'", query)

}
