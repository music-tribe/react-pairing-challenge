package get

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
)

//go:generate mockgen -destination=./mocks/get.go -package=getmocks -source=get.go
type GetDatabase interface {
	Get(userId, taskId uuid.UUID) (*domain.Task, error)
}

type GetRequest struct {
	UserId uuid.UUID `param:"userId" validate:"required"`
	TaskId uuid.UUID `param:"taskId" validate:"required"`
}

func Get(db GetDatabase) func(echo.Context) error {
	if db == nil {
		panic("get.Get: db has nil value")
	}

	return func(c echo.Context) error {
		req := GetRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.New().Struct(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		task, err := db.Get(req.UserId, req.TaskId)
		if err != nil {
			if err == database.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, task)
	}
}
