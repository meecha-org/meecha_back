package main

import (
	"net/http"
	"new-meecha/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	router := echo.New()

	// router.Use(middleware.Recover())
	router.Use(middleware.Logger())
	router.Use(middlewares.PocketAuth())

	router.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK,echo.Map{
			"result" : "hello world",
		})
	})

	router.Logger.Fatal(router.Start(":8090"))
}
