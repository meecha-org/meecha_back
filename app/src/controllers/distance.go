package controllers

import (
	"net/http"
	"new-meecha/logger"
	"new-meecha/services"
	"new-meecha/utils"

	"github.com/labstack/echo/v4"
)

//設定距離取得
func GetDistance(ctx echo.Context) error {
	// ユーザー情報を取得 (送信者)
	myid := ctx.Get("UserID").(string)

	// 通知距離を取得
	distance,err := services.GetDistance(myid)
	logger.Println("aaa",distance,err)
	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.JSON(http.StatusOK,distance)
}

type DistanceArgs struct {
	Distance int64 `json:"Distance"`
}

func UpdateDistance(ctx echo.Context) error {
	// ユーザー情報を取得 (送信者)
	myid := ctx.Get("UserID").(string)

	var args DistanceArgs
	// リクエストをバインド
	if err := ctx.Bind(&args); err != nil {
		utils.Println("failed to bind json : " + err.Error())
		return ctx.NoContent(http.StatusBadRequest)
	}

	// リクエストを取得
	err := services.UpdateDistance(myid,args.Distance)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.NoContent(http.StatusOK)
}
