package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/music-tribe/react-pairing-challenge/database"
	"github.com/music-tribe/react-pairing-challenge/handlers/add"
	"github.com/music-tribe/react-pairing-challenge/handlers/get"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())

	db, err := database.OpenMongoConnection("mongodb://root:rootpassword@localhost:27017/admin", e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm Alive!!!")
	})

	e.POST("/api/add", add.Add(db))
	e.GET("/api/get/:taskId", get.Get(db))

	e.Logger.Fatal(e.Start(":8083"))
}
