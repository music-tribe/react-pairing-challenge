package database

import (
	"context"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
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

func TestMongoDatabase_Add(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test as too long")
	}

	url := "mongodb://root:rootpassword@localhost:27017/admin"
	logger := echo.New().Logger
	db, err := OpenMongoConnection(url, logger)
	assert.NoError(t, err)
	defer db.CloseMongoConnection()

	t.Run("When the record already exists, we should get an error", func(t *testing.T) {
		task := domain.Task{
			Id:          uuid.New(),
			Name:        "done",
			Description: "exists",
		}

		err := db.Add(&task)
		assert.NoError(t, err)
		err = db.Add(&task)
		assert.ErrorIs(t, err, ErrDuplicate)
	})

	t.Run("When we add a task, we can retrieve it", func(t *testing.T) {
		expect := domain.Task{
			Id:          uuid.New(),
			Name:        "done",
			Description: "exists",
		}

		err := db.Add(&expect)
		assert.NoError(t, err)

		actual, err := db.Get(expect.Id)
		assert.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})
}
