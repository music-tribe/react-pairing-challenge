package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/music-tribe/react-pairing-challenge/database"
	_ "github.com/music-tribe/react-pairing-challenge/docs/features-api"
	"github.com/music-tribe/react-pairing-challenge/handlers/add"
	"github.com/music-tribe/react-pairing-challenge/handlers/delete"
	"github.com/music-tribe/react-pairing-challenge/handlers/get"
	"github.com/music-tribe/react-pairing-challenge/handlers/getall"
	"github.com/music-tribe/react-pairing-challenge/handlers/update"
	"github.com/music-tribe/react-pairing-challenge/handlers/upvote"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title features API
// @version 1.0
// @description This API allows you to create, get, alter and delete features.

// @contact.name MCloud Team
// @contact.email cloud@musictribe.com

// @host localhost:8083
// @BasePath /
// @schemes http
func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	db, err := database.OpenMongoConnection(os.Getenv("DB_URL"), e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/status", Status)

	grp := e.Group("/api")
	grp.POST("/:userId", add.Add(db))
	grp.GET("/:userId", getall.GetAll(db))
	grp.PUT("/:userId", update.Update(db))
	grp.GET("/:userId/:featureId", get.Get(db))
	grp.DELETE("/:userId/:featureId", delete.Delete(db))

	grp.PUT("/vote/:featureId", upvote.Upvote(db))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8083"))
}

// Status godoc
// @Summary Show if the server is alive.
// @Description get the status of server.
// @Accept */*
// @Produce text/plain
// @Router /status [get]
// @Success 200 {string} string "I'm Alive!!!"
func Status(c echo.Context) error {
	return c.String(http.StatusOK, "I'm Alive!!!")
}
