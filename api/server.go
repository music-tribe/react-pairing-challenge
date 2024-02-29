package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/handlers/add"
	"github.com/music-tribe/react-pairing-challenge/handlers/get"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())

	db, err := database.OpenMongoConnection(os.Getenv("DB_URL"), e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm Alive!!!")
	})

	e.POST("/api/add/:userId", add.Add(db))
	e.GET("/api/get/:userId/:taskId", get.Get(db))

	e.Logger.Fatal(e.Start(":8083"))
}
