package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"sqt/command"
	"sqt/config"
	"sqt/dataAdapter/redis"
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

	config.ReadClientConfigs()
	redis.Init()

	addr += ":" + config.Values.ConnPort

	start := time.Now()

	go runQuery(addr, key, start, hlth)
	go runQuery(addr, key, start, hlth)
	go runQuery(addr, key, start, hlth)

	time.Sleep(10 * time.Second)
}

func runQuery(address string, key string, start time.Time, hlth bool) {
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
		return
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting client at "+address+":", err.Error())
		return
	}

	_, err = fmt.Fprintf(conn, commString)
	if err != nil {
		fmt.Println("Error sending command to client at "+address+":", err.Error())
		return
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		fmt.Println("Error reading response from client at "+address+":", err.Error())
		return
	}

	result, err := message.Deserialize(response)
	if err != nil {
		fmt.Println("Error during deserializing response from "+address+":", err.Error())
		return
	}

	localVal, err := runQueryLocal(key)
	if err != nil {
		if err.Error() == ERROR_KEY_NOT_FOUND_REDIS {
			fmt.Println(ERROR_KEY_NOT_FOUND_REDIS)
			if result.Status == message.STATUS_OK_DB || result.Status == message.STATUS_OK_REDIS {
				fmt.Println("Saving received value to local storage instead")
				saveKeyValueToLocal(key, result.Data)
			}
		} else if err.Error() == ERROR_NO_CONNECTION_REDIS {
			fmt.Println(ERROR_NO_CONNECTION_REDIS)
			fmt.Println("Omitting compare part...")
		} else {
			fmt.Println("Unhandled error during connection to redis:", err.Error())
			fmt.Println("Omitting compare part...")
		}
	} else {
		fmt.Println("****************")
		if localVal == result.Data {
			fmt.Println("Received value (" + result.Data + ") is equal to locally stored value!")
		} else {
			fmt.Println("Received value (" + result.Data + ") is not equal to locally stored value (" + localVal + ")!")
		}
		fmt.Println("****************")
	}

	_ = localVal

	fmt.Println("\nTimestamp: " + strconv.Itoa(int(time.Since(start).Milliseconds())))
	fmt.Println("Status: " + message.STATUSES_TEXTS[result.Status])
	fmt.Println("Data: " + result.Data)
	fmt.Println("Time elapsed on query: " + strconv.Itoa(result.TimeElapsed))
	fmt.Println("Time elapsed total (query + possible queue): " + strconv.Itoa(result.TimeElapsedTotal))
	fmt.Println("Queue size at the time when request was received: " + strconv.Itoa(result.QueueSize))
}

func runQueryLocal(key string) (string, error) {
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

func printClientHelp() {
	fmt.Println("Execute me in form \"client {key} {host}\".")
}
