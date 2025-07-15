package controllers

import (
	"net/http"
	"new-meecha/logger"
	"new-meecha/models"
	"new-meecha/services"
	"new-meecha/utils"

	"github.com/labstack/echo/v4"
)


func UpdateIgnores(ctx echo.Context) error {
	// ユーザー情報を取得
	// myid := ctx.Get("UserID").(string)
	myid := ctx.Request().Header.Get("UserID")

	// bind
	var args []models.IgnoresArgs
	if err := ctx.Bind(&args); err != nil {
		utils.Println("failed to bind json : " + err.Error())
		return ctx.NoContent(http.StatusBadRequest)
	}

	// 通知しない距離を変更
	err := services.UpdateIgnores(myid,args)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.NoContent(http.StatusOK)
}


func GetIgnores(ctx echo.Context) error {
	// ユーザー情報を取得
	// myid := ctx.Get("UserID").(string)
	myid := ctx.Request().Header.Get("UserID")
	logger.Println(myid)
	// 通知しない距離を変更
	result,err := services.GetIgnores(myid)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.JSON(http.StatusOK,result)
}
