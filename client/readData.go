package main

import (
	"errors"
	"fmt"
	"net"
)

var packetHandlers map[int]func([]byte, *int) = map[int]func([]byte, *int){
	0: readEcho,
	1: readWelcome,
	2: readChatMessage,
}

func readDataFromServer(conn net.Conn) error {
	// Create a buffer to store the client's input
	buffer := make([]byte, 4096)

	// Read the input from the client
	unreadLength, err := conn.Read(buffer)
	if unreadLength == 0 || err != nil {
		return errors.New("connection closed")
	}

	// Store the received input from the server
	packet := buffer[:unreadLength]
	readPos := 0

	packetLength, packetID := unwrapPacket(packet, &readPos)
	if packetLength == -1 {
		return errors.New("error reading packet length")
	}

	value, isValid := packetHandlers[packetID]
	if !isValid {
		return errors.New("invalid packet id")
	}

	value(packet, &readPos)
	return nil
}

func readWelcome(packet []byte, readPos *int) {
	msg := readString(packet, readPos)
	id := readInt(packet, readPos)
	myID = id

	fmt.Printf("New message from server: \"%s\"\n", msg)
	fmt.Printf("Your assigned id is: \"%d\"\n", id)

	sendWelcomeReceived(username, myID)
}

func readEcho(packet []byte, readPos *int) {
	fmt.Println("Echo received.")
}

func readChatMessage(packet []byte, readPos *int) {
	messageType := readInt(packet, readPos)
	sender := readString(packet, readPos)
	msg := readString(packet, readPos)

	fmt.Printf("Received chat message \"%s\" of type \"%d\" from client %s.\n", msg, messageType, sender)
}
