package models

import (
)

type Distance struct {
	UserId    string `gorm:"primaryKey;column:user_id"`
	Distance  int64  //通知距離 (メートル)
	CreatedAt int64  //作成時間 (UnixTime)
	UpdatedAt int64  //更新時間 (UnixTime)
}

func GetDistance(uid string) (int64, error) {
	request := Distance{}

	// データベースから取得
	result := dbconn.Where(&Distance{
		UserId: uid,
	}).First(&request)

	// エラー処理
	if result.Error != nil {
		return request.Distance, result.Error
	}

	return request.Distance, nil
}
