package models

import "time"

type Ignores struct {
	Uid       string      //送信者ID
	Latitude  float64 
	Longitude float64 
	Size      int64  
	CreateAt  int64
}

type IgnoresArgs struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	Size      int64   `json:"Size"`
}


// uidに関連する除外ポイントを削除
func RemoveIgnores(uid string) error {
	return dbconn.Where(&Ignores{
		Uid:	uid,
	}).Unscoped().Delete(&Ignores{}).Error
}


// uidに関連する除外ポイントを登録
func SaveIgnores(uid string,ignores []IgnoresArgs) error {
	var ignoresRecords []Ignores

	for _, ignores := range ignores {
		ignoresRecords = append(ignoresRecords, Ignores{
			Uid:       uid,
			Latitude:  ignores.Latitude,
			Longitude: ignores.Longitude,
			Size:      ignores.Size,
			CreateAt:  time.Now().Unix(),
		})
	}

	return dbconn.Create(&ignoresRecords).Error
}

// uidに関連する除外ポイントを削除
func GetIgnores(uid string) ([]Ignores,error) {

	//除外ポイント
	ignores := []Ignores{}

	// リクエスト取得
	result := dbconn.Where(&Ignores{
		Uid:	uid,
	}).Find(&ignores)

		// エラー処理
	if result.Error != nil {
		return []Ignores{}, result.Error
	}

	return ignores,nil
}

