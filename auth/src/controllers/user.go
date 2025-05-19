package controllers

import (
	"auth/logger"
	"auth/models"
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetMe(ctx echo.Context) error {
	// セッションを取得
	session, ok := ctx.Get("session").(*models.Session)

	// エラー処理
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	// 自身の情報を取得
	user, err := services.GetMe(session.UserID)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}

func UpdateUser(ctx echo.Context) error {
	// リクエストボディを取得
	args := services.UpdateUserData{}

	// bind する
	if err := ctx.Bind(&args); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// ユーザーを更新する
	if err := services.UpdateUser(args); err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}	

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}

// ユーザーを削除する
func DeleteOauth(ctx echo.Context) error {
	// ヘッダからID取得
	userid := ctx.Request().Header.Get("userid")

	// 削除する
	err := services.DeleteUser(userid)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError,echo.Map{
			"error" : err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK,echo.Map{
		"result" : "success",
	})
}


func GetAllUsers(ctx echo.Context) (error) {
	// サービスを呼び出す
	users, err := services.GetUsers()

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, users)
}

// BAN を切り替える
func ToggleBan(ctx echo.Context) error {
	// bind する
	banArgs := services.BanArgs{}

	if err := ctx.Bind(&banArgs); err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusBadRequest,echo.Map{
			"result" : err.Error(),
		})
	}

	// BAN を切り替える
	err := services.ToggleBan(banArgs)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError,echo.Map{
			"result" : err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK,echo.Map{
		"result" : "success",
	})
}