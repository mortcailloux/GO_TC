package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func client(portString string) {
	// Connect to the server
	conn, err := net.Dial("tcp", portString)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to the server.")

	// Reader for server responses
	reader := bufio.NewReader(conn)
	// Reader for user input
	userInputReader := bufio.NewReader(os.Stdin)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Server closed the connection.")
				break
			}
			fmt.Println("Error reading from server:", err)
			break
		}

		fmt.Print("Message from server: " + message)

		if message == "Running program...\n" {
			fmt.Println("Waiting for program output...")
			// Step 2: Receive the program output from the server
			for {
				programOutput, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						fmt.Println("Server closed the connection.")
						break
					}
					fmt.Println("Error reading program output:", err)
					break
				}

				// Print the program output
				fmt.Print("Program output: " + programOutput)

				// Exit if the program indicates it's done
				if programOutput == "Program finished.\n" {
					fmt.Println("Interaction completed. Closing connection.")
					return
				}
			}
			break
		}

		fmt.Print(">> ")
		response, err := userInputReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		_, err = io.WriteString(conn, response)
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}
