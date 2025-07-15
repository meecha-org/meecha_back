package services

import "new-meecha/models"

//現在の距離を取得
func GetDistance(uid string) (int64,error){
	//設定した距離を取得
	distance,err := models.GetDistance(uid)
	if err != nil {
		return 0,err
	}
	return distance,nil
}

//現在の通知距離の変更
func UpdateDistance(uid string,distance int64)error{
	//設定した距離に変更
	err := models.UpdateDistance(uid,distance)
	if err != nil {
		return err
	}
	return nil
}