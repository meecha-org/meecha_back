package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()
	router.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, World!")
	})
	router.Logger.Fatal(router.Start(":8090"))
}