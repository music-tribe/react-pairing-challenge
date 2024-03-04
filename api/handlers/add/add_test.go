package add

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
	addmocks "github.com/music-tribe/react-pairing-challenge/handlers/add/mocks"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	e := echo.New()

	t.Run("when the db has a nil value, we should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Add(nil)
		})
	})

	t.Run("when the userId is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := addmocks.NewMockAddDatabase(ctrl)

		byt := []byte(`{"name":"hello","description":"do something","completed":false}`)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues("")

		err := Add(db)(ctx)
		assert.ErrorContains(t, err, "invalid UUID length: 0,")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when the feature name is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := addmocks.NewMockAddDatabase(ctrl)

		byt := []byte(`{"name":"","description":"do something","completed":false}`)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(uuid.New().String())

		err := Add(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'Name' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when the feature description is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := addmocks.NewMockAddDatabase(ctrl)

		byt := []byte(`{"name":"hello","description":"","completed":false}`)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(uuid.New().String())

		err := Add(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'Description' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, getStatusCode(rec, err))
	})

	t.Run("when we attempt to enter a duplicate record we should return a 409 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := addmocks.NewMockAddDatabase(ctrl)

		id := uuid.New()
		userId := uuid.New()
		byt := []byte(`{"name":"hello","description":"do something", "id":"` + id.String() + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(userId.String())

		db.EXPECT().Add(&domain.Feature{
			Id:          id,
			UserId:      userId,
			Name:        "hello",
			Description: "do something",
			Completed:   false,
		}).Return(database.ErrDuplicate)

		err := Add(db)(ctx)
		assert.ErrorContains(t, err, database.ErrDuplicate.Error())
		assert.Equal(t, http.StatusConflict, getStatusCode(rec, err))
	})

	t.Run("when we get an unknown error from the db we should return a 500 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := addmocks.NewMockAddDatabase(ctrl)

		id := uuid.New()
		userId := uuid.New()
		byt := []byte(`{"name":"hello","description":"do something", "id":"` + id.String() + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(userId.String())

		db.EXPECT().Add(&domain.Feature{
			Id:          id,
			UserId:      userId,
			Name:        "hello",
			Description: "do something",
			Completed:   false,
		}).Return(errors.New("some error"))

		err := Add(db)(ctx)
		assert.ErrorContains(t, err, "some error")
		assert.Equal(t, http.StatusInternalServerError, getStatusCode(rec, err))
	})

	t.Run("when the request is healthy, we should receive no errors and a 201 response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := addmocks.NewMockAddDatabase(ctrl)
		id := uuid.New()
		userId := uuid.New()

		byt := []byte(`{"name":"hello","description":"do something", "id":"` + id.String() + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(byt))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(userId.String())

		db.EXPECT().Add(&domain.Feature{
			Id:          id,
			UserId:      userId,
			Name:        "hello",
			Description: "do something",
			Completed:   false,
		}).Return(nil)

		err := Add(db)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, getStatusCode(rec, err))

		ar := new(AddResponse)
		err = json.Unmarshal(rec.Body.Bytes(), ar)
		assert.NoError(t, err)
		assert.Equal(t, id, ar.Id)
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
