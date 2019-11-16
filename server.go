package main

import (
	"fmt"
	"net"
	"os"
	"sqt/command"
	"sqt/config"
	"sqt/dataAdapter/redis"
	"sqt/message"
	"sqt/queue"
	"time"
)

var tasksCount = 0

func main() {
	config.ReadServerConfigs()
	redis.Init()

	l, err := net.Listen(config.Values.ConnType, ":"+config.Values.ConnPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Listening on port " + config.Values.ConnPort)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	unixTimeStart := int(time.Now().Unix())

	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading connection:", err.Error())
		return
	}

	commandToParse := string(buf[:reqLen])
	receivedCommand, err := command.Deserialize(commandToParse)
	if err != nil {
		fmt.Println("Error during reserializing incoming command:", err.Error())
		return
	}

	var resultMessage message.Message

	if receivedCommand.Type != command.COMMAND_RUN_QUEUE && receivedCommand.Type != command.COMMAND_HEALTHCHECK {
		fmt.Println(message.STATUSES_TEXTS[message.STATUS_WRONG_COMMAND_TYPE])

		unixTimeEnd := int(time.Now().Unix())
		resultMessage = message.Message{
			IsExecuted: false,
			TimeStart:  unixTimeStart,
			TimeEnd:    unixTimeEnd,
			Status:     message.STATUS_WRONG_COMMAND_TYPE,
			Command:    receivedCommand.Type,
			Key:        receivedCommand.KeyToCheck,
		}
	} else if len(receivedCommand.KeyToCheck) == 0 {
		fmt.Println(message.STATUSES_TEXTS[message.STATUS_WRONG_COMMAND_KEY])

		unixTimeEnd := int(time.Now().Unix())
		resultMessage = message.Message{
			IsExecuted: false,
			TimeStart:  unixTimeStart,
			TimeEnd:    unixTimeEnd,
			Status:     message.STATUS_WRONG_COMMAND_KEY,
			Command:    receivedCommand.Type,
			Key:        receivedCommand.KeyToCheck,
		}
	} else {
		fmt.Println("Received query for key " + receivedCommand.KeyToCheck)

		queueChannel := make(chan message.Message)

		commandExecutionMode := queue.MODE_QUEUE

		if receivedCommand.Type == command.COMMAND_HEALTHCHECK {
			fmt.Println("Healthcheck - query is executing immediately")
			commandExecutionMode = queue.MODE_HEALTHCHECK
		} else {
			tasksCount++
		}
		go queue.Run(receivedCommand.KeyToCheck, tasksCount, commandExecutionMode, queueChannel)
		resultMessage = <-queueChannel

		unixTimeEnd := int(time.Now().Unix())
		resultMessage.TimeStart = unixTimeStart
		resultMessage.TimeEnd = unixTimeEnd

		fmt.Println("Query for key "+receivedCommand.KeyToCheck,
			"- result status: \""+message.STATUSES_TEXTS[resultMessage.Status]+"\",",
			"result data: \""+resultMessage.Data+"\"")

		if receivedCommand.Type == command.COMMAND_HEALTHCHECK {

		} else {
			tasksCount--
		}
	}

	valueToReturn, err := message.Serialize(resultMessage)
	if err != nil {
		fmt.Println("Message serializing error:", err.Error())
	}

	test, err2 := message.Deserialize(valueToReturn)
	if err2 != nil {
		fmt.Println("Error during reserializing serialized message:", err2.Error())
	}
	_ = test

	_, _ = conn.Write([]byte(valueToReturn))
}
