package models

import (
	"new-meecha/utils"
	"time"
)

type Ignores struct {
	IgnoreId  string  `gorm:"primaryKey"`   //範囲ID
	Uid       string    					//ユーザーID
	Latitude  float64 						//緯度
	Longitude float64 						//経度
	Size      int64  						//円のサイズ
	CreateAt  int64							//作られた時間
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

		//ignoreId生成
		ignoreId,err := utils.Genid()
		if err != nil {
			return err
		}

		ignoresRecords = append(ignoresRecords, Ignores{
			IgnoreId:  ignoreId,
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

