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

type AddRequest struct {
	Id          uuid.UUID   `json:"-" bson:"_id"`
	UserId      uuid.UUID   `json:"userId" param:"userId" bson:"userId" validate:"required" example:"202c25c4-b2ce-4514-9045-890a1aa896ea"`
	Name        string      `json:"name" validate:"required" example:"My New Feature Request"`
	Description string      `json:"description" validate:"required" example:"Could we have this new feature please?"`
	Votes       []uuid.UUID `json:"votes" example:"['155dccaa-0299-4018-ab6b-90b9ee448943','ef2a27c4-b03d-4190-86f2-b1dc2538243e']"`
}

type AddResponse struct {
	Id uuid.UUID `json:"id"`
}

// Add godoc
// @Summary Add a new feature for this user.
// @Description Add a new feature for this user id.
// @Accept application/json
// @Produce application/json
// @Param feature body AddRequest true "Feature"
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
