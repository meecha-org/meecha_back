package services

import (
	"errors"
	"new-meecha/models"
	rediscache "new-meecha/redis-cache"
	"new-meecha/utils"
	"slices"
)

type IgnoresArgs struct {
	Size      int64   `json:"size"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// 除外ポイントの作成、更新
func UpdateIgnores(myid string, args []IgnoresArgs) error {
	for _, arg := range args {
		if !slices.Contains(ValidationList, arg.Size) {
			return errors.New("InvalidDistance")
		}
	}

	//myidに関連する除外ポイントを削除
	err := models.RemoveIgnores(myid)
	if err != nil {
		return err
	}

	// データを加工する
	addDatas := []models.IgnoresArgs{}

	// 引数を回す
	for _, arg := range args {
		// UUIDを生成する
		uid,_ := utils.Genid()

		//送られてきたポイントを登録する
		addDatas = append(addDatas, models.IgnoresArgs{
			Size:      arg.Size,
			Latitude:  arg.Latitude,
			Longitude: arg.Longitude,
			IgnoreId:  uid,
		})
	}

	if args != nil {
		//送られてきたポイントを登録する
		err = models.SaveIgnores(myid, addDatas)
		if err != nil {
			return err
		}
	}

	// キャッシュに追加できるように変更
	IgnoreAddDatas := []rediscache.IgnorePoint{}

	for _, arg := range addDatas {
		// データを整形して追加
		IgnoreAddDatas = append(IgnoreAddDatas, rediscache.IgnorePoint{
			Size:      arg.Size,
			Latitude:  arg.Latitude,
			Longitude: arg.Longitude,
			PointId:   arg.IgnoreId,
		})
	}

	// キャッシュを更新する
	err = rediscache.AddCacheIgnores(rediscache.CacheIgnoreArgs{
		UserID: myid,
		Datas:  IgnoreAddDatas,
	})

	return err
}

// 除外ポイントの作成、更新
func GetIgnores(myid string) ([]models.Ignores, error) {
	//myidに関連する除外ポイントを取得
	ignores, err := models.GetIgnores(myid)
	if err != nil {
		return []models.Ignores{}, err
	}

	return ignores, nil
}
