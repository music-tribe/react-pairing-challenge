package getall

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/domain"
	getAllmocks "github.com/music-tribe/react-pairing-challenge/handlers/getAll/mocks"
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

	t.Run("when we getAll an unknown error from the db we should return a 500 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getAllmocks.NewMockGetAllDatabase(ctrl)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		db.EXPECT().GetAll().Return([]*domain.Task{}, errors.New("some error"))

		err := GetAll(db)(ctx)
		assert.ErrorContains(t, err, "some error")
		assert.Equal(t, http.StatusInternalServerError, getAllStatusCode(rec, err))
	})

	t.Run("when the request is well formed, we should getAll the task back and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db := getAllmocks.NewMockGetAllDatabase(ctrl)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectTasks := []*domain.Task{
			{
				Id:          uuid.New(),
				Name:        "one",
				Description: "is it done yet",
			},
			{
				Id:          uuid.New(),
				Name:        "two",
				Description: "is it done yet",
			},
			{
				Id:          uuid.New(),
				Name:        "three",
				Description: "is it done yet",
			},
		}

		db.EXPECT().GetAll().Return(expectTasks, nil)

		err := GetAll(db)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, getAllStatusCode(rec, err))

		actualTasks := make([]*domain.Task, 0)
		err = json.Unmarshal(rec.Body.Bytes(), &actualTasks)
		assert.NoError(t, err)
		assert.Equal(t, expectTasks, actualTasks)
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
