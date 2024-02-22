package database

import (
	"context"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestOpenMongoSession(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test as too long")
	}

	t.Run("open the mongo connection", func(t *testing.T) {
		url := "mongodb://root:rootpassword@localhost:27017/admin"
		logger := echo.New().Logger
		db, err := OpenMongoConnection(url, logger)
		defer db.CloseMongoConnection()

		assert.NoError(t, err)
		err = db.client.Ping(context.TODO(), readpref.Nearest())
		assert.NoError(t, err)
	})
}
