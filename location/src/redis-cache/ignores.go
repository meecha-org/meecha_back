package rediscache

import (
	"context"
	"location/logger"

	// "new-meecha/utils"
	"github.com/redis/go-redis/v9"
)

type CacheIgnoreArgs struct {
	UserID string        //相手のユーザーID
	Datas  []IgnorePoint //除外ポイント一覧
}
type IgnorePoint struct {
	PointId   string  //除外ポイントID
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
		// ユーザーごとに redis に保存
		result := IgnoreRedisConn.GeoAdd(ctx, args.UserID, &redis.GeoLocation{
			Name:      point.PointId,
			Longitude: point.Longitude,
			Latitude:  point.Latitude,
		})

		// エラー処理
		if result.Err() != nil {
			return result.Err()
		}
	}

	return nil
}

// 近くの円を検索するエンドポイント
type SearchIgnoreArgs struct {
	UserID       string  //自分のユーザーID
	Latitude     float64 //緯度
	Longitude    float64 //経度
	SearchRadius int64   //検索する半径
}

type SearchIgnoreResult struct {
	Latitude  float64 //緯度
	Longitude float64 //経度
	Dist      float64 //距離
	IgnoreId  string  //除外ポイントID
}

func SearchCacheIgnores(args SearchIgnoreArgs) ([]SearchIgnoreResult, error) {
	// ユーザーIDをもとに検索する
	ignores, err := IgnoreRedisConn.GeoRadius(context.Background(), args.UserID, args.Longitude, args.Latitude, &redis.GeoRadiusQuery{
		Radius:    float64(args.SearchRadius),
		Unit:      "m",
		WithCoord: true,
		WithDist:  true,
		Sort:      "ASC",
	}).Result()

	// エラー処理
	if err != nil {
		return nil, err
	}

	logger.Println("redis ignores", ignores)

	// データを整形して返す
	returnData := make([]SearchIgnoreResult, len(ignores))

	for i, ignore := range ignores {
		// 検索結果を返す
		returnData[i] = SearchIgnoreResult{
			Latitude:  ignore.Latitude,
			Longitude: ignore.Longitude,
			Dist:      ignore.Dist,
			IgnoreId:  ignore.Name,
		}
	}

	return returnData, nil
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
