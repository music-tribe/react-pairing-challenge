package delete

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/uuid"
)

//go:generate mockgen -destination=./mocks/delete.go -package=deletemocks -source=delete.go
type DeleteDatabase interface {
	Delete(userId, taksId uuid.UUID) error
}

type DeleteRequest struct {
	UserId uuid.UUID `param:"userId" validate:"required"`
	TaskId uuid.UUID `param:"taskId" validate:"required"`
}

// Delete godoc
// @Summary Delete a users task.
// @Description Delete one of this users tasks.
// @Accept application/json
// @Produce text/plain
// @Param userId path string true "User UUID"
// @Param taskId path string true "Task UUID"
// @Router /api/{userId}/{taskId} [delete]
// @Success 200 {string} string "DELETED"
// @failure 400 {object} error
// @failure 404 {object} error
// @failure 500 {object} error
func Delete(db DeleteDatabase) func(echo.Context) error {
	if db == nil {
		panic("delete.Delete: db has nil value")
	}

	return func(c echo.Context) error {
		req := DeleteRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.New().Struct(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		err := db.Delete(req.UserId, req.TaskId)
		if err != nil {
			if err == database.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.String(http.StatusOK, "DELETED")
	}
}
