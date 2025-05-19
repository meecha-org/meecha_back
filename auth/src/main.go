package main

import (
	"auth/grpc"
	"auth/models"
	"auth/oauth2"
	"auth/services"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func Init() {
	// モデル初期化
	models.Init()

	// サービス初期化
	services.Init()

	// goth 初期化
	oauth2.InitGothic()

	// GRPC 初期化
	go grpc.RunServer()
}

func main() {
	// 初期化
	Init()

	// エンジン初期化
	// 認証初期化
	oauth2.UseProviders()

	// ルータ
	router := echo.New()

	// ルーティング設定
	SetupRouter(router)

	// ヘルスチェック
	router.GET("/health", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})
	
	router.Logger.Fatal(router.Start(":8080"))
}
