package add

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
)

//go:generate mockgen -destination=./mocks/add.go -package=addmocks -source=add.go
type AddDatabase interface {
	Add(task *domain.Task) error
}

type AddResponse struct {
	Id uuid.UUID `json:"id"`
}

func Add(db AddDatabase) func(echo.Context) error {
	if db == nil {
		panic("add.Add: db has nil value")
	}

	return func(c echo.Context) error {
		task := domain.Task{}

		if err := c.Bind(&task); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if task.Id == uuid.Nil {
			task.Id = uuid.New()
		}

		if err := validator.New().Struct(&task); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := db.Add(&task); err != nil {
			if err == database.ErrDuplicate {
				return echo.NewHTTPError(http.StatusConflict, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, AddResponse{Id: task.Id})
	}
}
