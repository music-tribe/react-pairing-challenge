package database

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMongoDatabase_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test as too long")
	}

	url := "mongodb://root:rootpassword@localhost:27017/admin"
	logger := echo.New().Logger
	db, err := OpenMongoConnection(url, logger)
	assert.NoError(t, err)
	defer db.CloseMongoConnection()

	t.Run("When the record can't be found, we should get an error", func(t *testing.T) {
		err := db.Delete(uuid.New(), uuid.New())
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("When the task exists, we can delete it", func(t *testing.T) {
		expect := domain.Task{
			Id:          uuid.New(),
			UserId:      uuid.New(),
			Name:        "done",
			Description: "exists",
		}

		err := db.Add(&expect)
		assert.NoError(t, err)

		err = db.Delete(expect.UserId, expect.Id)
		assert.NoError(t, err)

		_, err = db.Get(expect.UserId, expect.Id)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
