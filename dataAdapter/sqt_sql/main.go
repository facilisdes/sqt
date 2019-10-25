package sqt_sql

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"sqt/config"
)

const (
	ERROR_KEY_NOT_FOUND  = "Mysql key not found"
	ERROR_SQL_NO_CONNECT = "Cannot establish a connection to a database"
)

func OpenConnection() (*sql.DB, error) {
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", config.Values.DbLogin+":"+config.Values.DbPassword+"@tcp("+config.Values.DbHost+":"+
		config.Values.DbPort+")/"+config.Values.DbTable)

	if err != nil {
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

func CloseConnection(db *sql.DB) {
	_ = db.Close()
}

func GetSqlValue(key string) (string, error) {
	db, err := OpenConnection()
	if err != nil {
		//return "", errors.New(ERROR_SQL_NO_CONNECT)
		return "", err
	}

	defer CloseConnection(db)

	if db == nil {
		return "", errors.New(ERROR_SQL_NO_CONNECT)
	}
	rows, err := db.Query("select " + config.Values.DbValueColumnName + " from " + config.Values.DbTable + " where " +
		config.Values.DbKeyColumnName + "='" + key + "'")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	value := ""

	isChanged := false
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return "", err
		}
		isChanged = true
	}
	if !isChanged {
		return "", errors.New(ERROR_KEY_NOT_FOUND)
	}
	return value, nil
}
