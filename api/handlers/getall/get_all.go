package getall

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/domain"
)

//go:generate mockgen -destination=./mocks/get_all.go -package=getallmocks -source=get_all.go
type GetAllDatabase interface {
	GetAll() ([]*domain.Task, error)
}

func GetAll(db GetAllDatabase) func(echo.Context) error {
	if db == nil {
		panic("getAll.GetAll: db has nil value")
	}

	return func(c echo.Context) error {
		tasks, err := db.GetAll()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, tasks)
	}
}
