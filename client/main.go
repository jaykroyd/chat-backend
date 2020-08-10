package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var address string = "79.125.17.207"
var port string = "8080"
var myID int = -1
var connection net.Conn = nil
var username string = ""

func main() {
	fmt.Println("Enter ip address:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	if scanner.Text() != "" {
		address = scanner.Text()
	}

	for {
		fmt.Println("Enter username:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if scanner.Text() == "" {
			fmt.Println("Invalid username!")
		} else {
			username = scanner.Text()
			break
		}
	}

	socket := address + ":" + port
	fmt.Println("Starting client on socket: " + socket)

	// Attempt to connect to the server
	if err := connect(); err != nil {
		fmt.Println("Fatal main error")
		log.Fatal(err)
	}
}

func connect() error {
	socket := address + ":" + port
	// Connect to the server address via tcp protocol
	conn, err := net.Dial("tcp", socket)
	if err != nil {
		return err
	}
	fmt.Println("Successfully connected to", socket)
	connection = conn

	// Defer the closure of the connection, to ensure resources are released once the loop is finished
	defer conn.Close()

	go listenToServer(conn)
	checkForClientInput(conn)

	return nil
}

func checkForClientInput(conn net.Conn) {
	// Create scanner to read input in real time
	scanner := bufio.NewScanner(os.Stdin)

	// Initial message
	fmt.Println("Please type in your input.")

	// Start loop to listen for input from client
	for scanner.Scan() {
		fmt.Printf("Sending input to server: \"%s\"\n", scanner.Text())

		// Write the message to the server
		sendChatMessage(scanner.Text(), 1)
	}

	// Return any errors identified by the scanner
	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error")
	}
}

func listenToServer(conn net.Conn) {
	// Start loop to listen for new input from existing connection
	for {
		// Read the reply from the server
		if err := readDataFromServer(conn); err != nil && err != io.EOF {
			fmt.Println("Error reading packets")
		} else if err == io.EOF {
			fmt.Println("disconnected from server")
			connection = nil
			return
		}
	}
}
