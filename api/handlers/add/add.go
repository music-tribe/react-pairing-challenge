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
	Add(feature *domain.Feature) error
}

type AddResponse struct {
	Id uuid.UUID `json:"id"`
}

type Error echo.HTTPError

// Add godoc
// @Summary Add a new feature for this user.
// @Description Add a new feature for this user id.
// @Accept application/json
// @Produce application/json
// @Param feature body domain.Feature true "Feature"
// @Param userId path string true "User UUID"
// @Router /api/{userId} [post]
// @Success 200 {object} AddResponse
// @failure 400 {object} error
// @failure 409 {object} error
// @failure 500 {object} error
func Add(db AddDatabase) func(echo.Context) error {
	if db == nil {
		panic("add.Add: db has nil value")
	}

	return func(c echo.Context) error {
		feature := domain.Feature{}

		if err := c.Bind(&feature); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if feature.Id == uuid.Nil {
			feature.Id = uuid.New()
		}

		if err := validator.New().Struct(&feature); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := db.Add(&feature); err != nil {
			if err == database.ErrDuplicate {
				return echo.NewHTTPError(http.StatusConflict, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, AddResponse{Id: feature.Id})
	}
}
