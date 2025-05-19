package main

import (
	"location/controllers"
	"location/middlewares"
	redisfriend "location/redis-friend"
	"location/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// redis 初期化
	redisfriend.Init()

	// サービス初期化
	services.Init()

	// ミドルウェア初期化
	middlewares.Init()
	
	// ルーター
	router := echo.New()

	// router.Use(middleware.Recover())
	router.Use(middleware.Logger())
	// router.Use(middlewares.PocketAuth())

	router.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"result": "hello world",
		})
	}, middlewares.RequireAuth)

	// 位置情報をポストする関数
	router.POST("/update",controllers.UpdateLocation,middlewares.RequireAuth)

	router.Logger.Fatal(router.Start(":8090"))
}