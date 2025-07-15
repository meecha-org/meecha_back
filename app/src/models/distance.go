package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
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

//通知距離更新
func UpdateDistance(uid string, distance int64) error {
		dist := Distance{
		UserId:   uid,
		Distance: distance,
		UpdatedAt: time.Now().Unix(),
	}

	// 存在するか確認
	var existing Distance
	err := dbconn.Where("user_id = ?", uid).First(&existing).Error

	if err == nil {
		// 存在する場合は更新
		return dbconn.Model(&existing).Updates(Distance{
			Distance: distance,
			UpdatedAt: time.Now().Unix(),
		}).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 存在しない場合は新規作成
		dist.CreatedAt = time.Now().Unix() // CreatedAt を設定
		return dbconn.Create(&dist).Error
	}

	// その他のエラーが発生した場合
	return err
}
