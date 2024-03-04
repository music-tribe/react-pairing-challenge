package upvote

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
)

var (
	errVotedForOwnFeature = errors.New("Upvote: this user is attempting to vote for their own feature")
	errVoteAlreadyCounted = errors.New("Upvote: this user has already voted for the feature")
)

//go:generate mockgen -destination=./mocks/upvote.go -package=upvotemocks -source=upvote.go
type UpvoteDatabase interface {
	GetById(featureId uuid.UUID) (*domain.Feature, error)
	Update(feature *domain.Feature) error
}

type UpvoteRequest struct {
	UserId    uuid.UUID `json:"userId" form:"userId" validate:"required" example:"202c25c4-b2ce-4514-9045-890a1aa896ea"`
	FeatureId uuid.UUID `json:"-" param:"featureId" validate:"required" example:"b1f01569-ecff-4c60-a716-435b2e51f1ff"`
}

type UpvoteResponse struct {
	FeatureId uuid.UUID `json:"featureId" form:"featureId" example:"b1f01569-ecff-4c60-a716-435b2e51f1ff"`
	VoteCount int64     `json:"voteCount" example:"45928"`
}

// Upvote godoc
// @Summary Enables the user to vote for a new feature request.
// @Description Enables the user to place one vote against another users feature request.
// @Accept application/json
// @Produce text/plain
// @Param featureId path string true "Feature ID"
// @Param upvoteRequest body UpvoteRequest true "Upvote Request Body"
// @Router /api/vote/{featureId} [put]
// @Success 200 {object} UpvoteResponse
// @failure 400 {object} error
// @failure 404 {object} error
// @failure 409 {object} error
// @failure 500 {object} error
func Upvote(db UpvoteDatabase) func(echo.Context) error {
	if db == nil {
		panic("update.Upvote: db has nil value")
	}

	return func(c echo.Context) error {
		req := UpvoteRequest{}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.New().Struct(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		feature, err := db.GetById(req.FeatureId)
		if err != nil {
			if err == database.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if feature.UserId == req.UserId {
			return echo.NewHTTPError(http.StatusBadRequest, errVotedForOwnFeature)
		}

		for _, user := range feature.Votes {
			if user == req.UserId {
				return echo.NewHTTPError(http.StatusConflict, errVoteAlreadyCounted)
			}
		}

		feature.Votes = append(feature.Votes, req.UserId)

		if err := db.Update(feature); err != nil {
			if err == database.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, UpvoteResponse{
			FeatureId: req.FeatureId,
			VoteCount: int64(len(feature.Votes)),
		})
	}
}
