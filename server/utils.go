package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

func indexOfClient(element net.Conn, data []client) int {
	for k, v := range data {
		if element == v.conn {
			return k
		}
	}

	// Not Found
	return -1
}

func indexOfServer(element *Server) int {
	for k, v := range activeServerList {
		if element == v {
			return k
		}
	}

	// Not Found
	return -1
}

func addConnectedClient(c client, s *Server) error {
	if slot := c.id; s.connectedPlayersList[slot].conn != nil {
		return errors.New("Player already exists in slot " + strconv.Itoa(slot))
	}

	s.connectedPlayersList[c.id] = c
	fmt.Println("New connection accepted:", c.username)
	println("Connected Players: " + strconv.Itoa(getConnectedPlayersAmount(s)))

	return nil
}

func removeConnectedClient(conn net.Conn, s *Server) {
	i := indexOfClient(conn, s.connectedPlayersList)
	s.connectedPlayersList[i] = client{}
	fmt.Printf("Client %s was disconnected\n", conn.RemoteAddr())
	println("Connected Players: " + strconv.Itoa(getConnectedPlayersAmount(s)))
}

func assignServerSlot(s *Server) int {
	for i, r := range s.connectedPlayersList {
		if r.conn == nil {
			return i
		}
	}
	return -1
}

func getConnectedPlayersAmount(s *Server) int {
	c := 0
	for _, r := range s.connectedPlayersList {
		if r.conn != nil {
			c++
		}
	}
	return c
}

func setUsername(clientID int, username string) {
	activeServerList[0].connectedPlayersList[clientID].username = username
}
