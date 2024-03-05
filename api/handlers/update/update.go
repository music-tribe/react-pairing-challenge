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
	Update(feature *domain.Feature) error
}

// Update godoc
// @Summary Get all of a users features.
// @Description Get a all features releted to this userId.
// @Accept application/json
// @Produce text/plain
// @Param userId path string true "User UUID"
// @Param feature body domain.Feature true "Feature"
// @Router /api/{userId} [put]
// @Success 200 {object} domain.Feature
// @failure 400 {object} error
// @failure 404 {object} error
// @failure 500 {object} error
func Update(db UpdateDatabase) func(echo.Context) error {
	if db == nil {
		panic("update.Update: db has nil value")
	}

	return func(c echo.Context) error {
		feature := domain.Feature{}

		if err := c.Bind(&feature); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.New().Struct(&feature); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := db.Update(&feature); err != nil {
			if err == database.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, feature)
	}
}
