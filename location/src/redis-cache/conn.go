package rediscache

import (
	"location/logger"
	"location/models"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	friendConn      *redis.Client = nil
	IgnoreRedisConn *redis.Client = nil
)

func getRedis(db int) *redis.Client {
	if os.Getenv("REDIS_TYPE") == "redis" {
		// 通常のredisの場合
		// redis に接続
		return redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
			PoolSize: 1000,
		})
	}

	if os.Getenv("REDIS_TYPE") == "sentinel" {
		// sentinel の場合
		// sentinel に接続
		return redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs:     []string{os.Getenv("REDIS_HOST")},
			Password: os.Getenv("REDIS_PASSWORD"),
			MasterName: os.Getenv("REDIS_MASTER_NAME"),
			DB:       db,
			PoolSize: 1000,
		})
	}

	// panic 起こす
	panic("redis type error")
}

func Init() {
	// モデル初期化
	models.Init()

	// グローバル変数に格納
	friendConn = getRedis(0)

	// redis に接続
	IgnoreConn := getRedis(1)

	// グローバル変数に格納
	IgnoreRedisConn = IgnoreConn

	logger.Println("redis connected")
}
