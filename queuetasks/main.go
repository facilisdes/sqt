package queuetasks

import (
	"sqt/dataAdapter/redis"
	"sqt/dataAdapter/sqt_sql"
	"sqt/message"
)

func GetData(key string) (string, int) {
	value, err := getFromRedis(key)
	if err == nil {
		return value, message.STATUS_OK_REDIS
	}

	value, err = getFromSql(key)
	if err == nil {
		cacheKeyValuePair(key, value)
		return value, message.STATUS_OK_DB
	} else {
		if err.Error() == sqt_sql.ERROR_KEY_NOT_FOUND {
			return "", message.STATUS_ENTRY_NOT_FOUND
		} else {
			return "", message.STATUS_NO_ACTIVE_STORAGE
		}
	}
}

func getFromRedis(key string) (string, error) {
	value, err := redis.GetRedisValue(key)
	return value, err
}

func getFromSql(key string) (string, error) {
	value, err := sqt_sql.GetSqlValue(key)
	return value, err
}

func cacheKeyValuePair(key string, value string) {
	redis.SetRedisValue(key, value)
}
