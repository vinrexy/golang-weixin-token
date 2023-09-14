package utils

import "time"

// NowTimeStamp 当前时间戳

func NowTimeStamp() int64 {
	return time.Now().Unix()
}
