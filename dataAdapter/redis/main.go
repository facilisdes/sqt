package redis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"sqt/config"
	"time"
)

const (
	ERROR_KEY_NOT_FOUND    = "Redis key not found!"
	ERROR_REDIS_NO_CONNECT = "Cannot establish a connection to redis!"
)

var pool *redis.Pool

func Init() {
	// Initialize a connection pool and assign it to the pool global variable.
	pool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.Values.CacheHost+":"+config.Values.CachePort)
		},
	}
}

func GetRedisValue(key string) (string, error) {
	if pool == nil {
		return "", errors.New(ERROR_REDIS_NO_CONNECT)
	}
	conn := pool.Get()

	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))

	if err != nil {
		return "", err
	} else if len(value) == 0 {
		return "", errors.New(ERROR_KEY_NOT_FOUND)
	}

	return value, nil
}
func SetRedisValue(key string, value string) {
	if pool != nil {
		conn := pool.Get()
		defer conn.Close()
		conn.Do("SET", key, value)
	}
}
