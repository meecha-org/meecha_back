package controllers

import (
	"location/services"
	"location/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Location struct {
	// 緯度
	Lat float64 `json:"lat"`
	// 経度
	Lng float64 `json:"lng"`
}

func UpdateLocation(ctx echo.Context) error {
	// ユーザーを取得
	userid := ctx.Get("UserID").(string)

	// bind する
	location := Location{}
	if err := ctx.Bind(&location); err != nil {
		return err
	}

	// 現在時間を取得
	requestTime := utils.NowTime()

	// 位置情報を処理
	friends,err := services.UpdateLocation(services.Location{
		UserID:    userid,
		Latitude:  location.Lat,
		Longitude: location.Lng,
		Timestamp: requestTime,
	})

	// エラー処理
	if err != nil {
		utils.Println(err)
		return err
	}

	return ctx.JSON(http.StatusOK,friends)
}