package controllers

import (
	"net/http"
	"new-meecha/services"
	"new-meecha/utils"

	"github.com/labstack/echo/v4"
)

//設定距離取得
func GetDistance(ctx echo.Context) error {
	// ユーザー情報を取得 (送信者)
	myid := ctx.Get("UserID").(string)

	// リクエストを取得
	distance,err := services.GetDistance(myid)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.JSON(http.StatusOK,distance)
}