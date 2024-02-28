package get

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/domain"
	getmocks "github.com/music-tribe/react-pairing-challenge/handlers/get/mocks"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	e := echo.New()

	t.Run("when the db has a nil value, we should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Get(nil)
		})
	})

	t.Run("when the userId is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getmocks.NewMockGetDatabase(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "taskId")
		ctx.SetParamValues("", uuid.New().String())

		err := Get(db)(ctx)
		assert.ErrorContains(t, err, "invalid UUID length: 0")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when the userId has a nil value we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getmocks.NewMockGetDatabase(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "taskId")
		ctx.SetParamValues(uuid.Nil.String(), uuid.New().String())

		err := Get(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'UserId' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when the id is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getmocks.NewMockGetDatabase(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "taskId")
		ctx.SetParamValues(uuid.New().String(), "")

		err := Get(db)(ctx)
		assert.ErrorContains(t, err, "invalid UUID length: 0")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when the id has a nil value we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getmocks.NewMockGetDatabase(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "taskId")
		ctx.SetParamValues(uuid.New().String(), uuid.Nil.String())

		err := Get(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'TaskId' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when we the record can't be found we should return a 404 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getmocks.NewMockGetDatabase(ctrl)

		userId := uuid.New()
		id := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "taskId")
		ctx.SetParamValues(userId.String(), id.String())

		db.EXPECT().Get(userId, id).Return(nil, database.ErrNotFound)

		err := Get(db)(ctx)
		assert.ErrorContains(t, err, database.ErrNotFound.Error())
		assert.Equal(t, http.StatusNotFound, getStatusCode(rec, err))
	})

	t.Run("when we get an unknown error from the db we should return a 500 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getmocks.NewMockGetDatabase(ctrl)

		userId := uuid.New()
		id := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "taskId")
		ctx.SetParamValues(userId.String(), id.String())

		db.EXPECT().Get(userId, id).Return(nil, errors.New("some error"))

		err := Get(db)(ctx)
		assert.ErrorContains(t, err, "some error")
		assert.Equal(t, http.StatusInternalServerError, getStatusCode(rec, err))
	})

	t.Run("when the request is well formed, we should get the task back and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getmocks.NewMockGetDatabase(ctrl)

		userId := uuid.New()
		id := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "taskId")
		ctx.SetParamValues(userId.String(), id.String())

		expectTask := domain.Task{
			Id:          id,
			UserId:      userId,
			Name:        "blah",
			Description: "is it done yet",
		}

		db.EXPECT().Get(userId, id).Return(&expectTask, nil)

		err := Get(db)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, getStatusCode(rec, err))

		actualTask := new(domain.Task)
		err = json.Unmarshal(rec.Body.Bytes(), actualTask)
		assert.NoError(t, err)
		assert.Equal(t, expectTask, *actualTask)
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
