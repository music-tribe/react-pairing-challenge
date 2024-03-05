package getall

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
)

//go:generate mockgen -destination=./mocks/get_all.go -package=getallmocks -source=get_all.go
type GetAllDatabase interface {
	GetAll(userId uuid.UUID) ([]*domain.Feature, error)
}

type GetAllRequest struct {
	UserId uuid.UUID `param:"userId" validate:"required" example:"ef2a27c4-b03d-4190-86f2-b1dc2538243e"`
}

// GetAll godoc
// @Summary Get all of a users features.
// @Description Get a all features releted to this userId.
// @Accept application/json
// @Produce text/plain
// @Param userId path string true "User UUID"
// @Router /api/{userId} [get]
// @Success 200 {object} []domain.Feature
// @failure 400 {object} error
// @failure 404 {object} error
// @failure 500 {object} error
func GetAll(db GetAllDatabase) func(echo.Context) error {
	if db == nil {
		panic("getAll.GetAll: db has nil value")
	}

	return func(c echo.Context) error {
		req := GetAllRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.New().Struct(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		features, err := db.GetAll(req.UserId)
		if err != nil {
			if err == database.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, features)
	}
}
