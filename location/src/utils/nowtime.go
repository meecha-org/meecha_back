package utils

import "time"

func NowTime() int64 {
	return time.Now().Unix()
}