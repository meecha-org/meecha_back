package rediscache

import (
	"context"
	"new-meecha/logger"

	// "new-meecha/utils"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CacheIgnoreArgs struct {
	UserID string        //相手のユーザーID
	Datas  []IgnorePoint //除外ポイント一覧
}
type IgnorePoint struct {
	Size      int64   //半径
	Latitude  float64 //緯度
	Longitude float64 //経度
}

// キャッシュを追加
func AddCacheIgnores(args CacheIgnoreArgs) error {
	// バックグラウンドコンテキスト
	ctx := context.Background()

	// 全てのキャッシュを削除する
	err := IgnoreRedisConn.Del(ctx, args.UserID).Err()

	// エラー処理
	if err != nil {
		// キャッシュの削除に失敗時はエラーのみ出す
		logger.PrintErr(err)
	}

	// 除外ポイントのデータを回す
	for _, point := range args.Datas {
		// uuidを生成する
		uuid_obj, err := uuid.NewRandom()

		//エラー処理
		if err != nil {
			return err
		}

		// ユーザーごとに redis に保存
		result := IgnoreRedisConn.GeoAdd(ctx, args.UserID, &redis.GeoLocation{
			Name:      uuid_obj.String(),
			Longitude: point.Latitude,
			Latitude:  point.Longitude,
		})

		// エラー処理
		if result.Err() != nil {
			return result.Err()
		}
	}

	return nil
}

// キャッシュに存在するか
func ExistCacheIgnores(args CacheIgnoreArgs) bool {
	// バックグラウンドコンテキスト
	ctx := context.Background()

	// redis に保存
	count, err := IgnoreRedisConn.Exists(ctx, args.UserID).Result()

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
