package sqt_sql

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"sqt/config"
	"sqt/message"
)

const (
	ERROR_KEY_NOT_FOUND  = "Mysql key not found"
	ERROR_SQL_NO_CONNECT = "Cannot establish a connection to a database"
)

func OpenConnection() (*sql.DB, error) {
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", config.Values.DbLogin+":"+config.Values.DbPassword+"@tcp("+config.Values.DbHost+":"+
		config.Values.DbPort+")/"+config.Values.DbName)

	if err != nil {
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

func OpenConnectionToDB(database string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", config.Values.DbLogin+":"+config.Values.DbPassword+"@tcp("+config.Values.DbHost+":"+
		config.Values.DbPort+")/"+database)

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
func SaveEventData(msg message.Message, client string, localValue string) (int, error) {
	db, err := OpenConnectionToDB(config.MYSQL_SERVICE_DB)
	if err != nil {
		return -1, err
	}

	defer CloseConnection(db)

	if db == nil {
		return -1, errors.New(ERROR_SQL_NO_CONNECT)
	}

	ValueIsValidated := len(localValue) > 0 && msg.Data == localValue

	stmt, err := db.Prepare("INSERT " + config.MYSQL_EVENTS_TABLE + " SET IsExecuted=?,Status=?,StatusText=?," +
		"Data=?,LocalData=?,ValueIsValidated=?,TimeElapsed=?,TimeQueuedMin=?,TimeElapsedTotal=?,QueueSize=?,Command=?," +
		"RequestedKey=?,Client=?,TimeStart=?,TimeEnd=?")
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(msg.IsExecuted, msg.Status, message.STATUSES_TEXTS[msg.Status], msg.Data, localValue, ValueIsValidated, msg.TimeElapsed,
		msg.TimeQueuedMin, msg.TimeElapsedTotal, msg.QueueSize, msg.Command, msg.Key, client, msg.TimeStart, msg.TimeEnd)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}
