package tools

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	mqantutils "github.com/liangdas/mqant/utils"
)

var (
	RedisUrl = "redis://:@127.0.0.1:6379"

	tokenLockScript = `
		local lockKey = KEYS[1]
		local timeout = KEYS[2]
		local ret = redis.call("GET", lockKey)
		local locked = 0
		if not ret then
			redis.call("SET", lockKey, 1)
			redis.call("EXPIRE", lockKey, timeout)
			locked = 1
		end
		return {locked}
	`

	getTokenScript = `
		local tokenKey = KEYS[1]
		local token = redis.call("GET", tokenKey)
		if not token then
			return nil
		end
		local ttl = redis.call("TTL", tokenKey)
		return {token, ttl}
	`
)

func executeLua(conn redis.Conn, luaScript string, keyList ...interface{}) ([]interface{}, error) {
	script := redis.NewScript(len(keyList), luaScript)
	r, err := redis.Values(script.Do(conn, keyList...))
	if err != nil {
		return nil, err
	}
	if len(r) == 1 {
		v, err := redis.Int64(r[0], err)
		return []interface{}{v}, err
	} else {
		v1, err := redis.String(r[0], err)
		v2, err := redis.Int64(r[1], err)
		return []interface{}{v1, v2}, err
	}
}

func SetAccessToken(key string, value string, expire int64) (int64, error) {
	conn := mqantutils.GetRedisFactory().GetPool(RedisUrl).Get()
	defer conn.Close()
	if conn != nil {
		if expire > 0 {
			return redis.Int64(conn.Do("SET", key, value, "EX", expire))
		} else {
			return redis.Int64(conn.Do("SET", key, value))
		}
	}
	return 0, fmt.Errorf("redis connection failed")
}

func GetCacheData(key string) (string, error) {
	conn := mqantutils.GetRedisFactory().GetPool(RedisUrl).Get()
	defer conn.Close()
	if conn != nil {
		return redis.String(conn.Do("GET", key))
	}
	return "", fmt.Errorf("redis connection failed")
}

func GetCacheTTLData(key string) (int64, error) {
	conn := mqantutils.GetRedisFactory().GetPool(RedisUrl).Get()
	defer conn.Close()
	if conn != nil {
		return redis.Int64(conn.Do("TTL", key))
	}
	return 0, fmt.Errorf("redis connection failed")
}

func GetAccessToken(key string) ([]interface{}, error) {
	conn := mqantutils.GetRedisFactory().GetPool(RedisUrl).Get()
	defer conn.Close()
	if conn != nil {
		return executeLua(conn, getTokenScript, key)
	}
	return nil, fmt.Errorf("redis connection failed")
}

func LockToken(key string) ([]interface{}, error) {
	conn := mqantutils.GetRedisFactory().GetPool(RedisUrl).Get()
	defer conn.Close()
	if conn != nil {
		return executeLua(conn, tokenLockScript, key, 5)
	}
	return nil, fmt.Errorf("redis connection failed")
}

func DeleteKey(key string) {
	conn := mqantutils.GetRedisFactory().GetPool(RedisUrl).Get()
	defer conn.Close()
	if conn != nil {
		conn.Do("DEL", key)
	}
}

func PublishChannel(channel string, msg string) (int64, error) {
	conn := mqantutils.GetRedisFactory().GetPool(RedisUrl).Get()
	defer conn.Close()
	if conn != nil {
		return redis.Int64(conn.Do("PUBLISH", channel, msg))
	}
	return 0, fmt.Errorf("redis connection failed")
}

func ListenPubSubChannel(channel string) (string, error) {
	conn := mqantutils.GetRedisFactory().GetPool(RedisUrl).Get()
	defer conn.Close()
	if conn != nil {

		psc := redis.PubSubConn{Conn: conn}
		defer psc.Close()
		if err := psc.Subscribe(redis.Args{}.AddFlat(channel)...); err != nil {
			return "", err
		}

		// Start a goroutine to receive notifications from the server.
		for {
			switch n := psc.ReceiveWithTimeout(time.Second * 20).(type) {
			case error:
				return "", n
			case redis.Message:
				return string(n.Data), nil
			}
		}

		// Wait for goroutine to complete.
	}
	return "", fmt.Errorf("redis connection failed")
}
