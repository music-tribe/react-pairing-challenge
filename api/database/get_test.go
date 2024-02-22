package database

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMongoDatabase_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test as too long")
	}

	url := "mongodb://root:rootpassword@localhost:27017/admin"
	logger := echo.New().Logger
	db, err := OpenMongoConnection(url, logger)
	assert.NoError(t, err)
	defer db.CloseMongoConnection()

	t.Run("When the record can't be found, we should get an error", func(t *testing.T) {
		task, err := db.Get(uuid.New())
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, task)
	})

	t.Run("When the task exists, we can retrieve it", func(t *testing.T) {
		expect := domain.Task{
			Id:          uuid.New(),
			Name:        "done",
			Description: "exists",
		}

		t.Cleanup(func() {
			_ = db.Delete(expect.Id)
		})

		err := db.Add(&expect)
		assert.NoError(t, err)

		actual, err := db.Get(expect.Id)
		assert.NoError(t, err)
		assert.Equal(t, expect, *actual)
	})
}
