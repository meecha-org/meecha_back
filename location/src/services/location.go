package services

import (
	"context"
	redisfriend "location/redis-friend"
	"location/utils"
	"slices"

	"github.com/redis/go-redis/v9"
)

type Location struct {
	UserID    string  //	ユーザーID
	Latitude  float64 `json:"latitude"`  // 緯度
	Longitude float64 `json:"longitude"` // 経度
	Timestamp int64   `json:"timestamp"` // タイムスタンプ
}

type NearFriend struct {
	UserID    string  `json:"userid"`    //相手のユーザーID
	Latitude  float64 `json:"latitude"`  // 緯度
	Longitude float64 `json:"longitude"` // 経度
	Dist      float64 `json:"dist`       //距離
}

// 位置情報を更新して近くにいるフレンドを返す
func UpdateLocation(args Location) ([]NearFriend, error) {
	// redis に保存
	err := conn.GeoAdd(context.Background(), geoKey, &redis.GeoLocation{
		Name:      args.UserID,
		Longitude: args.Longitude,
		Latitude:  args.Latitude,
	}).Err()

	// エラー処理
	if err != nil {
		utils.Println(err)
		return []NearFriend{}, err
	}

	// フレンド情報を取得
	cached, err := redisfriend.GetCacheFriend(args.UserID)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return []NearFriend{}, err
	}

	// 半径5キロ以内のユーザー取得
	nearUsers, err := conn.GeoRadius(context.Background(), geoKey, args.Longitude, args.Latitude, &redis.GeoRadiusQuery{
		Radius:    5,
		Unit:      "km",
		WithCoord: true,
		WithDist:  true,
	}).Result()

	// エラー処理
	if err != nil {
		utils.Println(err)
		return []NearFriend{}, err
	}

	// 返すフレンド一覧
	retuurnFriends := []NearFriend{}

	for _, user := range nearUsers {
		// 近くにいるユーザーIDをとる
		targetId := user.Name

		// 自分のIDの場合戻る
		if targetId == args.UserID {
			continue
		}

		// フレンドかどうか
		if !slices.Contains(cached.FriendIds, targetId) {
			// フレンドではない場合
			continue
		}

		// 返すリストに追加
		retuurnFriends = append(retuurnFriends, NearFriend{
			UserID:    user.Name,
			Latitude:  user.Latitude,
			Longitude: user.Longitude,
			Dist:      user.Dist,
		})
	}

	return retuurnFriends, nil
}
