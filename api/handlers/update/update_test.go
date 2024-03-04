package update

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/domain"
	updatemocks "github.com/music-tribe/react-pairing-challenge/handlers/update/mocks"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	e := echo.New()

	t.Run("when the db has a nil value, we should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Update(nil)
		})
	})

	t.Run("when the feature name is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := updatemocks.NewMockUpdateDatabase(ctrl)

		byt := []byte(`{"name":"","description":"do something","completed":false}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		err := Update(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'Name' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, updateStatusCode(rec, err))
	})

	t.Run("when the feature description is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := updatemocks.NewMockUpdateDatabase(ctrl)

		byt := []byte(`{"name":"hello","description":"","completed":false}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		err := Update(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'Description' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, updateStatusCode(rec, err))
	})

	t.Run("when we the record can't be found we should return a 404 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := updatemocks.NewMockUpdateDatabase(ctrl)

		userId := uuid.New()
		id := uuid.New()
		byt := []byte(`{"name":"hello","description":"some description","completed":false, "id":"` + id.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(userId.String())

		expectfeature := &domain.Feature{
			Id:          id,
			UserId:      userId,
			Name:        "hello",
			Description: "some description",
		}

		db.EXPECT().Update(expectfeature).Return(database.ErrNotFound)

		err := Update(db)(ctx)
		assert.ErrorContains(t, err, database.ErrNotFound.Error())
		assert.Equal(t, http.StatusNotFound, updateStatusCode(rec, err))
	})

	t.Run("when we update an unknown error from the db we should return a 500 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := updatemocks.NewMockUpdateDatabase(ctrl)

		userId := uuid.New()
		id := uuid.New()
		byt := []byte(`{"name":"hello","description":"some description","completed":false, "id":"` + id.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(userId.String())

		expectfeature := &domain.Feature{
			Id:          id,
			UserId:      userId,
			Name:        "hello",
			Description: "some description",
		}

		db.EXPECT().Update(expectfeature).Return(errors.New("some error"))

		err := Update(db)(ctx)
		assert.ErrorContains(t, err, "some error")
		assert.Equal(t, http.StatusInternalServerError, updateStatusCode(rec, err))
	})

	t.Run("when the request is well formed we should return a 200 response and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := updatemocks.NewMockUpdateDatabase(ctrl)

		userId := uuid.New()
		id := uuid.New()
		byt := []byte(`{"name":"hello","description":"some description","completed":false, "id":"` + id.String() + `"}`)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(userId.String())

		expectfeature := &domain.Feature{
			Id:          id,
			UserId:      userId,
			Name:        "hello",
			Description: "some description",
		}

		db.EXPECT().Update(expectfeature).Return(nil)

		err := Update(db)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, updateStatusCode(rec, err))
	})
}

func updateStatusCode(rec *httptest.ResponseRecorder, err error) int {
	if err == nil {
		return rec.Code
	}

	hterr := &echo.HTTPError{}
	if errors.As(err, &hterr) {
		return hterr.Code
	}

	return 500
}
