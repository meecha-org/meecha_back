package redisfriend

import (
	"location/models"
	"location/utils"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	conn *redis.Client = nil
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
	conn = redisConn

	utils.Println("redis connected")
}
