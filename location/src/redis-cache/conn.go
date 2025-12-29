package rediscache

import (
	"context"
	"location/logger"
	"location/models"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConn struct {
	// redis
	Conn *redis.Client

	// redis cluster
	ClusterConn *redis.ClusterClient
}

func (rconn *RedisConn) Close() {
	if rconn.Conn != nil {
		rconn.Conn.Close()
	}
	if rconn.ClusterConn != nil {
		rconn.ClusterConn.Close()
	}
}

// set する
func (rconn *RedisConn) Set(ctx context.Context, key string, value []byte, expiry time.Duration) error {
	if rconn.Conn != nil {
		return rconn.Conn.Set(ctx, key, value, time.Second*time.Duration(expiry)).Err()
	} else {
		return rconn.ClusterConn.Set(ctx, key, value, time.Second*time.Duration(expiry)).Err()
	}
}

// get する
func (rconn *RedisConn) Get(ctx context.Context, key string) (string, error) {
	if rconn.Conn != nil {
		return rconn.Conn.Get(ctx, key).Result()
	} else {
		return rconn.ClusterConn.Get(ctx, key).Result()
	}
}

// del する
func (rconn *RedisConn) Del(ctx context.Context, key string) error {
	if rconn.Conn != nil {
		return rconn.Conn.Del(ctx, key).Err()
	} else {
		return rconn.ClusterConn.Del(ctx, key).Err()
	}
}

// exists する
func (rconn *RedisConn) Exists(ctx context.Context, key string) (int64, error) {
	if rconn.Conn != nil {
		return rconn.Conn.Exists(ctx, key).Result()
	} else {
		return rconn.ClusterConn.Exists(ctx, key).Result()
	}
}

// GeoAdd
func (rconn *RedisConn) GeoAdd(ctx context.Context, key string, location *redis.GeoLocation) (int64,error) {
	if rconn.Conn != nil {
		return rconn.Conn.GeoAdd(ctx, key, location).Result()
	} else {
		return rconn.ClusterConn.GeoAdd(ctx, key, location).Result()
	}
}

// GeoRadius
func (rconn *RedisConn) GeoRadius(ctx context.Context, key string,longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	if rconn.Conn != nil {
		return rconn.Conn.GeoRadius(ctx, key, longitude, latitude, query).Result()
	} else {
		return rconn.ClusterConn.GeoRadius(ctx, key, longitude, latitude, query).Result()
	}
}

// ZRem
func (rconn *RedisConn) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	if rconn.Conn != nil {
		return rconn.Conn.ZRem(ctx, key, members).Result()
	} else {
		return rconn.ClusterConn.ZRem(ctx, key, members).Result()
	}
}

var (
	friendConn      *RedisConn = nil
	IgnoreRedisConn *RedisConn = nil
)

func getRedis(db int) *RedisConn {
	if os.Getenv("REDIS_TYPE") == "redis" {
		// 通常のredisの場合
		// redis に接続
		rconn := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       db,
			PoolSize: 1000,
		})

		return &RedisConn{
			Conn: rconn,
		}
	}

	if os.Getenv("REDIS_TYPE") == "sentinel" {
		// sentinel の場合
		// sentinel に接続
		clusterConn := redis.NewFailoverClusterClient(&redis.FailoverOptions{
			MasterName:     os.Getenv("REDIS_MASTER_NAME"),
			SentinelAddrs:  []string{os.Getenv("REDIS_HOST")},
			Password:      	os.Getenv("REDIS_PASSWORD"),
			// DB:             db,
			PoolSize:       1000,
			RouteRandomly:  true,
			// 耐障害性のための設定
			MaxRetries:      3,
			MinRetryBackoff: time.Millisecond * 100,
			MaxRetryBackoff: time.Second * 2,
			// タイムアウト設定
			DialTimeout:  time.Second * 5,
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 3,
			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				err := cn.Select(ctx, db).Err()
				return err
			},
		})

		return &RedisConn{
			ClusterConn: clusterConn,
		}
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

