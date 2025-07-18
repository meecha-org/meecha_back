package rediscache

import (
	"location/models"
	"location/utils"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	friendConn      *redis.Client = nil
	IgnoreRedisConn *redis.Client = nil
)

func Init() {
	// モデル初期化
	models.Init()

	// redis に接続
	redisConn := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		PoolSize: 1000,
	})

	// グローバル変数に格納
	friendConn = redisConn

	// redis に接続
	IgnoreConn := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       1,
		PoolSize: 1000,
	})

	// グローバル変数に格納
	IgnoreRedisConn = IgnoreConn

	utils.Println("redis connected")
}
