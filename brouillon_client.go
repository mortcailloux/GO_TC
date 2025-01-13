package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func client(portString int) {
	// Connect to the server
	conn, err := net.Dial("tcp", portString)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to the server. Type your message:")

	// Read from stdin and send to server
	reader := bufio.NewReader(conn)

	for {
		fmt.Print(">> ")
		message := reader.ReadString('\n')

		// Send message to server
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		// Receive response from server
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Print("Server response: " + response)
	}
}
