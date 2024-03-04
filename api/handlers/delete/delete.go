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
	UserId    uuid.UUID `param:"userId" validate:"required"`
	FeatureId uuid.UUID `param:"featureId" validate:"required"`
}

// Delete godoc
// @Summary Delete a users feature.
// @Description Delete one of this users features.
// @Accept application/json
// @Produce text/plain
// @Param userId path string true "User UUID"
// @Param featureId path string true "Feature UUID"
// @Router /api/{userId}/{featureId} [delete]
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

		err := db.Delete(req.UserId, req.FeatureId)
		if err != nil {
			if err == database.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.String(http.StatusOK, "DELETED")
	}
}
