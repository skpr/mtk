package mock

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// GetDB which can be used for testing purposes.
func GetDB(t *testing.T) (*sql.Conn, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	conn, err := db.Conn(context.TODO())
	assert.Nil(t, err)
	return conn, mock
}
