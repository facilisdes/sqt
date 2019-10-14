package main

import (
	"fmt"
	"net"
	"os"
	"sqt/config"
	"sqt/message"
	"sqt/queue"
)

var tasksCount = 0

func main() {
	config.ReadConfigs()
	l, err := net.Listen(config.Values.ConnType, ":"+config.Values.ConnPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on port " + config.Values.ConnPort)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)

	}
}

func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	_ = reqLen
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	keyToRead := string(buf[:reqLen])

	fmt.Println("Received key:", keyToRead)


	queueChannel := make(chan message.Message)
	go queue.Run(keyToRead, tasksCount, queueChannel)
	tasksCount++
	resultMessage := <-queueChannel
	tasksCount--

	valueToReturn := message.ToGOB64(resultMessage)
	test := message.FromGOB64(valueToReturn)
	_=test;

	fmt.Println(valueToReturn)
	conn.Close()
}