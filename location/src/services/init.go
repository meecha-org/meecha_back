package services

import (
	"location/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	conn *redis.Client = nil
	cacheConn *redis.Client = nil
)

func Init() {
	// redis に接続
	redisConn := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       1,
		PoolSize: 1000,
	})

	// キャッシュ用のredis
	// redis に接続
	cacheRedis := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       2,
		PoolSize: 1000,
	})

	// グローバル変数に格納
	conn = redisConn
	cacheConn = cacheRedis

	utils.Println(conn)
	utils.Println(cacheConn)
	utils.Println("location redis connected")

	go func() {
		for {
			// 有効期限が切れたキーを削除する
			RemoveExpiryGeo()

			time.Sleep(time.Second * 3)
		}
	}()
}


