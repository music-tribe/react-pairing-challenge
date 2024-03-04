package upvote

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/domain"
	upvotemocks "github.com/music-tribe/react-pairing-challenge/handlers/upvote/mocks"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpvote(t *testing.T) {
	e := echo.New()

	t.Run("when the db has a nil value, we should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Upvote(nil)
		})
	})

	t.Run("when the userId is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		byt := []byte(`{"userId":"b1f01569-ecff-4c60-a716-435b2e51f1ff"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues("")

		err := Upvote(db)(ctx)
		assert.ErrorContains(t, err, "invalid UUID length: 0,")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when the userId has a nil value we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		byt := []byte(`{"userId":"b1f01569-ecff-4c60-a716-435b2e51f1ff"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues(uuid.Nil.String())

		err := Upvote(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'FeatureId' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when the featureId is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		byt := []byte(`{"userId":""}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues(uuid.New().String())

		err := Upvote(db)(ctx)
		assert.ErrorContains(t, err, "invalid UUID length: 0")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when the featureId has a nil value we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		byt := []byte(`{"userId":"` + uuid.Nil.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues(uuid.New().String())

		err := Upvote(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'UserId' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when we the record can't be found we should return a 404 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		userId := uuid.New()
		featureId := uuid.New()
		byt := []byte(`{"userId":"` + userId.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues(featureId.String())

		db.EXPECT().GetById(featureId).Return(nil, database.ErrNotFound)

		err := Upvote(db)(ctx)
		assert.ErrorContains(t, err, database.ErrNotFound.Error())
		assert.Equal(t, http.StatusNotFound, getStatusCode(rec, err))
	})

	t.Run("when there is an internal server error we should return a 500 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		userId := uuid.New()
		featureId := uuid.New()
		byt := []byte(`{"userId":"` + userId.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues(featureId.String())

		someError := errors.New("some error")
		db.EXPECT().GetById(featureId).Return(nil, someError)

		err := Upvote(db)(ctx)
		assert.ErrorContains(t, err, someError.Error())
		assert.Equal(t, http.StatusInternalServerError, getStatusCode(rec, err))
	})

	t.Run("when we the user attempts to vote for their own feature request, we should return a 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		userId := uuid.New()
		featureId := uuid.New()
		byt := []byte(`{"userId":"` + userId.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues(featureId.String())

		feature := domain.Feature{
			UserId:      userId,
			Id:          featureId,
			Name:        "something",
			Description: "hello thing",
			Votes: []uuid.UUID{
				userId, // a previous vote by this user
			},
		}

		db.EXPECT().GetById(featureId).Return(&feature, nil)

		err := Upvote(db)(ctx)
		assert.ErrorContains(t, err, errVotedForOwnFeature.Error())
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when we the user has already voted for this feature, we should return a 409", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		userId := uuid.New()
		featureId := uuid.New()
		byt := []byte(`{"userId":"` + userId.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues(featureId.String())

		feature := domain.Feature{
			UserId:      uuid.New(),
			Id:          featureId,
			Name:        "something",
			Description: "hello thing",
			Votes: []uuid.UUID{
				userId, // a previous vote by this user
			},
		}

		db.EXPECT().GetById(featureId).Return(&feature, nil)

		err := Upvote(db)(ctx)
		assert.ErrorContains(t, err, errVoteAlreadyCounted.Error())
		assert.Equal(t, http.StatusConflict, getStatusCode(rec, err))
	})

	t.Run("when we update database method errors we should return a 500 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		userId := uuid.New()
		featureId := uuid.New()
		byt := []byte(`{"userId":"` + userId.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues(featureId.String())

		feature := domain.Feature{
			UserId:      uuid.New(),
			Id:          featureId,
			Name:        "something",
			Description: "hello thing",
			Votes: []uuid.UUID{
				uuid.New(), // a previous vote
			},
		}

		db.EXPECT().GetById(featureId).Return(&feature, nil)

		someError := errors.New("some error")
		db.EXPECT().Update(&feature).Return(someError)

		err := Upvote(db)(ctx)
		assert.ErrorContains(t, err, someError.Error())
		assert.Equal(t, http.StatusInternalServerError, getStatusCode(rec, err))
	})

	t.Run("when the request is well formed we should return a 200 and an UpdateResponse object", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := upvotemocks.NewMockUpvoteDatabase(ctrl)

		userId := uuid.New()
		featureId := uuid.New()
		byt := []byte(`{"userId":"` + userId.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("featureId")
		ctx.SetParamValues(featureId.String())

		feature := domain.Feature{
			UserId:      uuid.New(),
			Id:          featureId,
			Name:        "something",
			Description: "hello thing",
			Votes: []uuid.UUID{
				uuid.New(), // a previous vote
			},
		}

		db.EXPECT().GetById(featureId).Return(&feature, nil)

		db.EXPECT().Update(&feature).Return(nil)

		expect := UpvoteResponse{
			FeatureId: featureId,
			VoteCount: int64(len(feature.Votes) + 1),
		}

		err := Upvote(db)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, getStatusCode(rec, err))

		actual := UpvoteResponse{}
		err = json.Unmarshal(rec.Body.Bytes(), &actual)
		assert.NoError(t, err)
		assert.Equal(t, expect, actual)
	})
}

func getStatusCode(rec *httptest.ResponseRecorder, err error) int {
	if err == nil {
		return rec.Code
	}

	hterr := &echo.HTTPError{}
	if errors.As(err, &hterr) {
		return hterr.Code
	}

	return 500
}
