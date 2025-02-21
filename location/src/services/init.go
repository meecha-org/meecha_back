package services

import (
	"location/utils"

	"github.com/redis/go-redis/v9"
)

var (
	conn *redis.Client = nil
	geoKey = "meecha_geo"
)

func Init() {
	// redis に接続
	redisConn := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       1,
		PoolSize: 1000,
	})

	// グローバル変数に格納
	conn = redisConn

	utils.Println("location redis connected")
}
