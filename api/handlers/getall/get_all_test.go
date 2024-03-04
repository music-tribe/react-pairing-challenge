package getall

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
	getallmocks "github.com/music-tribe/react-pairing-challenge/handlers/getall/mocks"
	"github.com/music-tribe/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	e := echo.New()

	t.Run("when the db has a nil value, we should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			GetAll(nil)
		})
	})

	t.Run("when the userId is missing we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getallmocks.NewMockGetAllDatabase(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues("")

		err := GetAll(db)(ctx)
		assert.ErrorContains(t, err, "message=invalid UUID length: 0")
		assert.Equal(t, http.StatusBadRequest, getAllStatusCode(rec, err))
	})

	t.Run("when the userId has a nil value we should return a 400 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getallmocks.NewMockGetAllDatabase(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(uuid.Nil.String())

		err := GetAll(db)(ctx)
		assert.ErrorContains(t, err, "Error:Field validation for 'UserId' failed on the 'required' tag")
		assert.Equal(t, http.StatusBadRequest, getAllStatusCode(rec, err))
	})

	t.Run("when the user can't be found in the db we should return a 404 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getallmocks.NewMockGetAllDatabase(ctrl)

		userId := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(userId.String())

		db.EXPECT().GetAll(userId).Return([]*domain.Feature{}, database.ErrNotFound)

		err := GetAll(db)(ctx)
		assert.ErrorContains(t, err, database.ErrNotFound.Error())
		assert.Equal(t, http.StatusNotFound, getAllStatusCode(rec, err))
	})

	t.Run("when we getAll an unknown error from the db we should return a 500 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getallmocks.NewMockGetAllDatabase(ctrl)

		userId := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(userId.String())

		db.EXPECT().GetAll(userId).Return([]*domain.Feature{}, errors.New("some error"))

		err := GetAll(db)(ctx)
		assert.ErrorContains(t, err, "some error")
		assert.Equal(t, http.StatusInternalServerError, getAllStatusCode(rec, err))
	})

	t.Run("when the request is well formed, we should getAll the feature back and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getallmocks.NewMockGetAllDatabase(ctrl)

		userId := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("userId")
		ctx.SetParamValues(userId.String())

		expectfeatures := []*domain.Feature{
			{
				Id:          uuid.New(),
				UserId:      userId,
				Name:        "one",
				Description: "is it done yet",
			},
			{
				Id:          uuid.New(),
				UserId:      userId,
				Name:        "two",
				Description: "is it done yet",
			},
			{
				Id:          uuid.New(),
				UserId:      userId,
				Name:        "three",
				Description: "is it done yet",
			},
		}

		db.EXPECT().GetAll(userId).Return(expectfeatures, nil)

		err := GetAll(db)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, getAllStatusCode(rec, err))

		actualfeatures := make([]*domain.Feature, 0)
		err = json.Unmarshal(rec.Body.Bytes(), &actualfeatures)
		assert.NoError(t, err)
		assert.Equal(t, expectfeatures, actualfeatures)
	})
}

func getAllStatusCode(rec *httptest.ResponseRecorder, err error) int {
	if err == nil {
		return rec.Code
	}

	hterr := &echo.HTTPError{}
	if errors.As(err, &hterr) {
		return hterr.Code
	}

	return 500
}
