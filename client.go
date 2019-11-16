package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"sqt/command"
	"sqt/config"
	"sqt/dataAdapter/redis"
	"sqt/dataAdapter/sqt_sql"
	"sqt/message"
	"strconv"
	"time"
)

const (
	ERROR_KEY_NOT_FOUND_REDIS = "Cannot compare received answer with local value - no key found"
	ERROR_NO_CONNECTION_REDIS = "Cannot compare received answer with local value - redis is unavailable"
)

func main() {
	var key string
	flag.StringVar(&key, "key", "42", "key to check")
	var addr string
	flag.StringVar(&addr, "host", "127.0.0.1", "host to get key from")
	var hlth bool
	flag.BoolVar(&hlth, "hc", false, "should we use healthcheck and execute command immediately? (default - false)")
	var commsCount int
	flag.IntVar(&commsCount, "c", 1, "number of requests to be send. 0 for infinite")
	var sendPeriodFrom int
	flag.IntVar(&sendPeriodFrom, "pf", 100, "minimal pause between requests")
	var sendPeriodTo int
	flag.IntVar(&sendPeriodTo, "pt", 5000, "maximal pause between requests")

	flag.Parse()
	if commsCount < 0 {
		fmt.Println("Wrong requests count - must be 0 for infinite or 1+ !")
		return
	}

	config.ReadClientConfigs()
	redis.Init()

	addr += ":" + config.Values.ConnPort

	runQueries(commsCount, key, sendPeriodFrom, sendPeriodTo, addr, hlth)
}

func runQueries(count int, key string, sendPeriodFrom int, sendPeriodTo int, address string, healthMode bool) {
	start := time.Now()
	runsCount := 0
	chans := make(chan bool)
	if count != 0 {
		chans = make(chan bool, count)
	}

	for count == 0 || runsCount < count {
		if healthMode == false {
			// if commsCount equals 0 or we haven't run required amount of commands
			var waitTime int
			if sendPeriodFrom == sendPeriodTo {
				waitTime = sendPeriodFrom
			} else {
				rand.Seed(time.Now().UnixNano())
				waitTime = sendPeriodFrom + rand.Intn(sendPeriodTo-sendPeriodFrom)
			}
			time.Sleep(time.Duration(waitTime) * time.Millisecond)
		}

		go runQueryRoutine(address, key, start, healthMode, chans)
		runsCount++
	}

	for i := 0; i < count; i++ {
		_ = <-chans
	}
}

func runQueryRoutine(address string, key string, start time.Time, hlth bool, queueChannel chan bool) {
	commType := command.COMMAND_RUN_QUEUE

	if hlth {
		commType = command.COMMAND_HEALTHCHECK
	}

	comm := command.Command{
		Type:       commType,
		KeyToCheck: key,
	}
	commString, err := command.Serialize(comm)
	commTest, _ := command.Deserialize(commString)
	_ = commTest
	if err != nil {
		fmt.Println("Error serializing command for "+address+":", err.Error())
		queueChannel <- true
		return
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting client at "+address+":", err.Error())
		queueChannel <- true
		return
	}

	_, err = fmt.Fprintf(conn, commString)
	if err != nil {
		fmt.Println("Error sending command to client at "+address+":", err.Error())
		queueChannel <- true
		return
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		fmt.Println("Error reading response from client at "+address+":", err.Error())
		queueChannel <- true
		return
	}

	result, err := message.Deserialize(response)
	if err != nil {
		fmt.Println("Error during deserializing response from "+address+":", err.Error())
		queueChannel <- true
		return
	}

	strToPrint := ""

	strToPrint += "\nTimestamp: " + strconv.Itoa(int(time.Since(start).Milliseconds())) + "\n"
	strToPrint += "Status: " + message.STATUSES_TEXTS[result.Status] + "\n"
	strToPrint += "Data: " + result.Data + "\n"

	localVal, err := runLocalQuery(key)
	if err != nil {
		if err.Error() == ERROR_KEY_NOT_FOUND_REDIS {
			strToPrint += ERROR_KEY_NOT_FOUND_REDIS + "\n"
			if result.Status == message.STATUS_OK_DB || result.Status == message.STATUS_OK_REDIS {
				strToPrint += "Saving received value to local storage instead\n"
				saveKeyValueToLocal(key, result.Data)
			}
		} else if err.Error() == ERROR_NO_CONNECTION_REDIS {
			strToPrint += ERROR_KEY_NOT_FOUND_REDIS + "\n"
			strToPrint += "Omitting compare part..." + "\n"
		} else {
			strToPrint += "Unhandled error during connection to redis:" + err.Error() + "\n"
			strToPrint += "Omitting compare part..." + "\n"
		}
	} else {
		if localVal == result.Data {
			strToPrint += "Received value (" + result.Data + ") is equal to locally stored value!" + "\n"
		} else {
			strToPrint += "Received value (" + result.Data + ") is not equal to locally stored value (" + localVal + ")!" + "\n"
		}
	}

	strToPrint += "Time elapsed on query: " + strconv.Itoa(result.TimeElapsed) + "\n"
	strToPrint += "Time elapsed total (query + possible queue): " + strconv.Itoa(result.TimeElapsedTotal) + "\n"
	strToPrint += "Queue size after request was received: " + strconv.Itoa(result.QueueSize) + "\n"

	i, err := sqt_sql.SaveEventData(result, address, localVal)
	if err != nil {
		strToPrint += "Error during saving data to a local storage:" + err.Error()
	} else {
		strToPrint += "Data saved to local storage! ID=" + strconv.Itoa(i)
	}
	fmt.Println(strToPrint)

	queueChannel <- true
}

func runLocalQuery(key string) (string, error) {
	value, err := redis.GetRedisValue(key)
	if err == nil {
		return value, nil
	}
	if err.Error() == redis.ERROR_KEY_NOT_FOUND {
		return "", errors.New(ERROR_KEY_NOT_FOUND_REDIS)
	}
	if err.Error() == redis.ERROR_REDIS_NO_CONNECT {
		return "", errors.New(ERROR_NO_CONNECTION_REDIS)
	}
	return "", err
}

func saveKeyValueToLocal(key string, value string) {
	redis.SetRedisValue(key, value)
}
