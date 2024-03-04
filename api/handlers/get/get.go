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
	Get(userId, featureId uuid.UUID) (*domain.Feature, error)
}

type GetRequest struct {
	UserId    uuid.UUID `param:"userId" validate:"required"`
	FeatureId uuid.UUID `param:"featureId" validate:"required"`
}

// Get godoc
// @Summary Get a users feature.
// @Description Get a feature with matching feature and user id.
// @Accept application/json
// @Produce text/plain
// @Param userId path string true "User UUID"
// @Param featureId path string true "Feature UUID"
// @Router /api/{userId}/{featureId} [get]
// @Success 200 {object} domain.Feature
// @failure 400 {object} error
// @failure 404 {object} error
// @failure 500 {object} error
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

		feature, err := db.Get(req.UserId, req.FeatureId)
		if err != nil {
			if err == database.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, feature)
	}
}
