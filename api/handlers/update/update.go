package update

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/domain"
)

//go:generate mockgen -destination=./mocks/update.go -package=updatemocks -source=update.go
type UpdateDatabase interface {
	Update(task *domain.Task) error
}

func Update(db UpdateDatabase) func(echo.Context) error {
	if db == nil {
		panic("update.Update: db has nil value")
	}

	return func(c echo.Context) error {
		task := domain.Task{}

		if err := c.Bind(&task); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.New().Struct(&task); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := db.Update(&task); err != nil {
			if err == database.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, task)
	}
}
