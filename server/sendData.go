package main

import (
	"fmt"
)

func sendDataToClient(clientID int, s *Server, packet *[]byte, packageID int) {
	writeID(packageID, packet)
	writeLength(packet)
	conn := s.connectedPlayersList[clientID].conn
	if _, err := conn.Write(*packet); err != nil {
		panic(err)
	}
}

func sendDataToAllClients(s *Server, packet *[]byte, packageID int) {
	writeID(packageID, packet)
	writeLength(packet)
	for _, r := range s.connectedPlayersList {
		if r.conn != nil {
			if _, err := r.conn.Write(*packet); err != nil {
				fmt.Println("Couldn't send packet to", r.addr)
			}
		}
	}
}

func sendDataToAllClientsExcept(exceptClientID int, s *Server, packet *[]byte, packageID int) {
	writeID(packageID, packet)
	writeLength(packet)
	for _, r := range s.connectedPlayersList {
		if r.conn != nil && r.id != exceptClientID {
			if _, err := r.conn.Write(*packet); err != nil {
				fmt.Println("Couldn't send packet to", r.addr)
			}
		}
	}
}

func sendWelcome(clientID int, s *Server, message string) {
	packet := make([]byte, 0)
	writeString(&packet, message)
	writeInt(&packet, clientID)

	sendDataToClient(clientID, s, &packet, 1)
	fmt.Println("Welcome packet sent to the client:", s.connectedPlayersList[clientID].addr)
}

func sendEcho(clientID int, s *Server, message string) {
	packet := make([]byte, 0)
	writeString(&packet, message)

	sendDataToClient(clientID, s, &packet, 0)
}

func sendChatMessage(clientID int, s *Server, message string, messageType int) {
	packet := make([]byte, 0)
	writeInt(&packet, messageType)
	sender := s.connectedPlayersList[clientID].username
	writeString(&packet, sender)
	writeString(&packet, message)

	sendDataToAllClientsExcept(clientID, s, &packet, 2)
}

func sendClientList(clientID int, s *Server) {
	packet := make([]byte, 0)

	// Get connected users
	usernames := make([]string, 0)
	ids := make([]int, 0)

	for _, r := range s.connectedPlayersList {
		if r.conn != nil {
			usernames = append(usernames, r.username)
			ids = append(ids, r.id)
		}
	}

	if len(usernames) != len(ids) {
		fmt.Println("number of client usernames and ids are different")
		return
	}

	// Write connected users to packet
	writeInt(&packet, len(usernames))

	for _, r := range usernames {
		writeString(&packet, r)
	}

	for _, r := range ids {
		writeInt(&packet, r)
	}

	sendDataToClient(clientID, s, &packet, 3)
}
