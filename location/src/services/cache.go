package services

import (
	"context"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

const (
	CacheExpiry = 300 //５分
)

type LocationCache struct {
	UserID          string   //ユーザーID
	LastNearFriends []string //最後に近くにいたフレンドID
	Timestamp       int64    //キャッシュ更新時間
}

// 位置情報の応答をキャッシュ
func SetCacheLocation(args LocationCache) error {
	// msgpack (マーシャリング)
	marshaled, err := msgpack.Marshal(args)

	// エラー処理
	if err != nil {
		return err
	}

	// キャッシュに保存
	err = cacheConn.Set(context.Background(), args.UserID, marshaled, time.Second*CacheExpiry).Err()

	return err
}

// 位置情報のキャッシュ取得
func GetCachedLocation(userid string) (LocationCache, error) {
	// キャッシュから取得
	data, err := cacheConn.Get(context.Background(), userid).Result()

	// エラー処理
	if err != nil {
		return LocationCache{}, err
	}

	// 戻した構造体
	unmarshaled := LocationCache{}

	// msgpack (アンマーシャル)
	err = msgpack.Unmarshal([]byte(data), &unmarshaled)

	// エラー処理
	if err != nil {
		return LocationCache{}, err
	}

	return unmarshaled, nil
}

// キャッシュに存在するか
func ExistLocationCache(userid string) bool {
	// バックグラウンドコンテキスト
	ctx := context.Background()

	// redis に保存
	count,err := cacheConn.Exists(ctx,userid).Result()

	// エラー処理
	if err != nil {
		return false
	}

	// 0以上なら存在する
	if count > 0 {
		return true
	}

	return false
}