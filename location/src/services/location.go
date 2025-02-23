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

type NearResponse struct {
	Removed     []string     `json:"removed"` //距離が離れたフレンド
	NearFriends []NearFriend `json:"near"`    //近くにいるフレンド
}

type NearFriend struct {
	UserID    string  `json:"userid"`    //相手のユーザーID
	Latitude  float64 `json:"latitude"`  // 緯度
	Longitude float64 `json:"longitude"` // 経度
	Dist      float64 `json:"dist`       //距離
}

// 位置情報を更新して近くにいるフレンドを返す
func UpdateLocation(args Location) (NearResponse, error) {
	// redis に保存
	err := conn.GeoAdd(context.Background(), geoKey, &redis.GeoLocation{
		Name:      args.UserID,
		Longitude: args.Longitude,
		Latitude:  args.Latitude,
	}).Err()

	// エラー処理
	if err != nil {
		utils.Println(err)
		return NearResponse{}, err
	}

	// フレンド情報を取得
	cached, err := redisfriend.GetCacheFriend(args.UserID)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return NearResponse{}, err
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
		return NearResponse{}, err
	}

	// 返すフレンド一覧
	retuurnFriends := []NearFriend{}
	NearFriendIds := []string{}

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

		// ID一覧の追加
		NearFriendIds = append(NearFriendIds, user.Name)
	}

	// キャッシュが存在するか
	if !ExistLocationCache(args.UserID) {
		// キャッシュが存在しない時

		// キャッシュを作成
		err := SetCacheLocation(LocationCache{
			UserID:          args.UserID,
			LastNearFriends: NearFriendIds,
			Timestamp:       utils.NowTime(),
		})

		// エラー処理
		if err != nil {
			utils.Println(err)
			return NearResponse{},err
		}
		
		// データを返す
		return NearResponse{
			Removed:     []string{},
			NearFriends: retuurnFriends,
		},nil
	}

	// キャッシュが存在する時
	// キャッシュを取得
	lastCached,err := GetCachedLocation(args.UserID)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return NearResponse{},err
	}

	// 範囲から外れたフレンドを取得
	nearRemoved := findRemovedElements(lastCached.LastNearFriends,NearFriendIds)

	// キャッシュを更新
	err = SetCacheLocation(LocationCache{
		UserID:          args.UserID,
		LastNearFriends: NearFriendIds,
		Timestamp:       utils.NowTime(),
	})

	return NearResponse{
		Removed:     nearRemoved,
		NearFriends: retuurnFriends,
	}, err
}

// 配列から消えた要素を取得する
func findRemovedElements(original []string, updated []string) []string {
	removed := []string{}
	updatedSet := make(map[string]struct{}, len(updated))

	// 更新された配列の要素をマップに格納
	for _, item := range updated {
		updatedSet[item] = struct{}{}
	}

	// 元の配列の要素をチェック
	for _, item := range original {
		if _, found := updatedSet[item]; !found {
			removed = append(removed, item)
		}
	}

	return removed
}
