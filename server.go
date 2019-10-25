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

		resultMessage = message.Message{
			IsExecuted: false,
			Status:     message.STATUS_WRONG_COMMAND_TYPE,
		}
	} else if len(receivedCommand.KeyToCheck) == 0 {
		fmt.Println(message.STATUSES_TEXTS[message.STATUS_WRONG_COMMAND_KEY])

		resultMessage = message.Message{
			IsExecuted: false,
			Status:     message.STATUS_WRONG_COMMAND_KEY,
		}
	} else {
		fmt.Println("Received query for key " + receivedCommand.KeyToCheck)

		queueChannel := make(chan message.Message)
		go queue.Run(receivedCommand.KeyToCheck, tasksCount, queueChannel)
		tasksCount++
		resultMessage = <-queueChannel

		fmt.Println("Query for key "+receivedCommand.KeyToCheck,
			"- result status: \""+message.STATUSES_TEXTS[resultMessage.Status]+"\",",
			"result data: \""+resultMessage.Data+"\"")

		tasksCount--
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
