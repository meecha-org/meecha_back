package main

import (
	"fmt"
	"location/controllers"
	"location/middlewares"
	redisfriend "location/redis-friend"
	"location/services"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// env 読み込み
	loadEnv()

	// redis 初期化
	redisfriend.Init()

	// サービス初期化
	services.Init()
	
	// ルーター
	router := echo.New()

	// router.Use(middleware.Recover())
	router.Use(middleware.Logger())
	// router.Use(middlewares.PocketAuth())

	router.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"result": "hello world",
		})
	}, middlewares.PocketAuth())

	// 位置情報をポストする関数
	router.POST("/update",controllers.UpdateLocation,middlewares.PocketAuth())

	router.Logger.Fatal(router.Start(":8090"))
}

// .envを呼び出します。
func loadEnv() {
	// ここで.envファイル全体を読み込みます。
	// この読み込み処理がないと、個々の環境変数が取得出来ません。
	// 読み込めなかったら err にエラーが入ります。
	err := godotenv.Load(".env")

	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}
