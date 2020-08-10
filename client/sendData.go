package main

import (
	"fmt"
)

func sendDataToServer(packet *[]byte, packageID int) {
	writeID(packageID, packet)
	writeLength(packet)
	if _, err := connection.Write(*packet); err != nil {
		panic(err)
	}
}

func sendWelcomeReceived(username string, clientID int) {
	packet := make([]byte, 0)
	writeString(&packet, username)
	writeInt(&packet, clientID)

	sendDataToServer(&packet, 1)
	fmt.Println("Welcome packet received confirmation sent to the server.")
}

func sendChatMessage(message string, messageType int) {
	packet := make([]byte, 0)
	writeInt(&packet, messageType)
	writeString(&packet, message)

	sendDataToServer(&packet, 2)
	fmt.Println("Chat message sent to the server.")
}
