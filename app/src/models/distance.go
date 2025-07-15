package models

import "time"

type Distance struct {
	UserId    string `gorm:"primaryKey"`
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

	// 存在するか確認し、存在すれば更新、存在しなければ挿入する
	return dbconn.Where("userid = ?", uid).Assign(dist).FirstOrCreate(&dist).Error
}
