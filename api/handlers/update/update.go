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

// Update godoc
// @Summary Get all of a users tasks.
// @Description Get a all tasks releted to this userId.
// @Accept application/json
// @Produce text/plain
// @Param userId path string true "User UUID"
// @Param task body domain.Task true "Task"
// @Router /api/{userId} [put]
// @Success 200 {object} domain.Task
// @failure 400 {object} error
// @failure 404 {object} error
// @failure 500 {object} error
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
