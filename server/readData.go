package main

import (
	"errors"
	"fmt"
	"net"
)

var packetHandlers map[int]func(client, *Server, []byte, *int) = map[int]func(client, *Server, []byte, *int){
	0:  readEcho,
	1:  readWelcomeReceived,
	2:  readChatMessage,
	99: readClientCommand,
}

func readDataFromClient(conn net.Conn, s *Server) error {
	// Create a buffer to store the client's input
	buffer := make([]byte, 4096)

	// Read the input from the client
	unreadLength, err := conn.Read(buffer)
	if unreadLength == 0 || err != nil {
		return errors.New("connection closed")
	}

	// Store the received input from the client
	client := s.connectedPlayersList[indexOfClient(conn, s.connectedPlayersList)]
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

	value(client, s, packet, &readPos)
	return nil
}

func readWelcomeReceived(c client, s *Server, packet []byte, readPos *int) {
	msg := readString(packet, readPos)
	id := readInt(packet, readPos)

	setUsername(id, msg)
	fmt.Printf("Welcome Message Received Confirmation from clientID %d: \"%s\"\n", id, msg)
}

func readEcho(c client, s *Server, packet []byte, readPos *int) {
	fmt.Println("Echo received.")
}

func readChatMessage(c client, s *Server, packet []byte, readPos *int) {
	messageType := readInt(packet, readPos)
	msg := readString(packet, readPos)

	sendChatMessage(c.id, s, msg, messageType)

	fmt.Printf("Received chat message \"%s\" of type \"%d\" from client %s.\n", msg, messageType, c.conn.RemoteAddr())
}

func readClientCommand(c client, s *Server, packet []byte, readPos *int) {
	commandID := readInt(packet, readPos)

	value, isValid := commandIndexList[commandID]
	if !isValid {
		fmt.Println("invalid command id")
		return
	}

	value(c, s)
}
