package services

import (
	"errors"
	"new-meecha/models"
	"slices"

	"gorm.io/gorm"
)

const (
	DefaultDistance = 3000
)

var ValidationList = []int64{50,200,500,1000,3000,5000}

//現在の距離を取得
func GetDistance(uid string) (int64,error){
	//設定した距離を取得
	distance,err := models.GetDistance(uid)
	
	//設定がなかった時
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return DefaultDistance,nil
	}

	if err != nil {
		return 0,err
	}
	return distance,nil
}

//現在の通知距離の変更
func UpdateDistance(uid string,distance int64)error{
	if !slices.Contains(ValidationList, distance) {
		return errors.New("InvalidDistance")
	}
	
	//設定した距離に変更
	err := models.UpdateDistance(uid,distance)
	if err != nil {
		return err
	}
	return nil
}