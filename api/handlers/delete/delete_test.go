package delete

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	deletemocks "github.com/music-tribe/react-pairing-challenge/handlers/delete/mocks"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	e := echo.New()

	t.Run("when the db has a nil value, we should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Delete(nil)
		})
	})

	t.Run("when the id is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := deletemocks.NewMockDeleteDatabase(ctrl)

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "featureId")
		ctx.SetParamValues(uuid.New().String(), "")

		err := Delete(db)(ctx)
		assert.ErrorContains(t, err, "invalid UUID length: 0")
		assert.Equal(t, http.StatusBadRequest, deleteStatusCode(rec, err))
	})

	t.Run("when the id has a nil value we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := deletemocks.NewMockDeleteDatabase(ctrl)

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "featureId")
		ctx.SetParamValues(uuid.New().String(), uuid.Nil.String())

		err := Delete(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'FeatureId' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, deleteStatusCode(rec, err))
	})

	t.Run("when we the record can't be found we should return a 404 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := deletemocks.NewMockDeleteDatabase(ctrl)

		userId := uuid.New()
		id := uuid.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "featureId")
		ctx.SetParamValues(userId.String(), id.String())

		db.EXPECT().Delete(userId, id).Return(database.ErrNotFound)

		err := Delete(db)(ctx)
		assert.ErrorContains(t, err, database.ErrNotFound.Error())
		assert.Equal(t, http.StatusNotFound, deleteStatusCode(rec, err))
	})

	t.Run("when we delete an unknown error from the db we should return a 500 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := deletemocks.NewMockDeleteDatabase(ctrl)

		userId := uuid.New()
		id := uuid.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "featureId")
		ctx.SetParamValues(userId.String(), id.String())

		db.EXPECT().Delete(userId, id).Return(errors.New("some error"))

		err := Delete(db)(ctx)
		assert.ErrorContains(t, err, "some error")
		assert.Equal(t, http.StatusInternalServerError, deleteStatusCode(rec, err))
	})

	t.Run("when the request is well formed, we should delete the feature back and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := deletemocks.NewMockDeleteDatabase(ctrl)

		userId := uuid.New()
		id := uuid.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId", "featureId")
		ctx.SetParamValues(userId.String(), id.String())

		db.EXPECT().Delete(userId, id).Return(nil)

		err := Delete(db)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, deleteStatusCode(rec, err))

		assert.Equal(t, "DELETED", rec.Body.String())
		assert.NoError(t, err)
	})
}

func deleteStatusCode(rec *httptest.ResponseRecorder, err error) int {
	if err == nil {
		return rec.Code
	}

	hterr := &echo.HTTPError{}
	if errors.As(err, &hterr) {
		return hterr.Code
	}

	return 500
}
