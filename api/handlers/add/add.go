package add

import (
	"github.com/labstack/echo/v4"
	"github.com/music-tribe/react-pairing-challenge/domain"
)

//go:generate mockgen -destination=./mocks/add.go -package=addmocks -source=add.go
type AddDatabase interface {
	Add(task *domain.Task) error
}

func Add(db AddDatabase) func(echo.Context) error {
	if db == nil {
		panic("add.Add: db has nil value")
	}

	return func(ctx echo.Context) error {
		return nil
	}
}
