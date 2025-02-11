package controllers

import (
	"net/http"
	"new-meecha/services"
	"new-meecha/utils"

	"github.com/labstack/echo/v4"
)

type SearchArgs struct {
	UserName string `json:"name"`
}

func SearchUser(ctx echo.Context) error {
	var args SearchArgs

	// リクエストをバインド
	if err := ctx.Bind(&args); err != nil {
		utils.Println("failed to bind json : " + err.Error())
		return ctx.NoContent(http.StatusBadRequest)
	}

	// バリデーション
	if args.UserName == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	// 検索する
	result,err := services.SearchByName(args.UserName)

	// エラー処理
	if err != nil {
		utils.Println("failed to find user : " + err.Error())
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK,result)
}