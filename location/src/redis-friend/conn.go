package redisfriend

import (
	"location/models"
	"location/utils"

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
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
		PoolSize: 1000,
	})

	// グローバル変数に格納
	conn = redisConn

	utils.Println("redis connected")
}
