package main

import (
	"fmt"
	"strconv"
	"strings"
)

var commandIndexList map[int]func(client, *Server) = map[int]func(client, *Server){
	0: requestClientList,
}

func checkForInput(input string, s *Server) {
	// Stop currently selected server
	if input == "quit" {
		s.Stop()
	}

	// Print all connected clients
	if input == "clients" {
		for _, r := range s.connectedPlayersList {
			fmt.Println("ID: "+strconv.Itoa(r.id)+" | Username: "+r.username+" | Address:", r.addr)
		}
	}

	// Print all active servers
	if input == "servers" {
		for _, r := range activeServerList {
			fmt.Println("ID: " + strconv.Itoa(indexOfServer(r)) + " | Address: " + r.listener.Addr().String())
		}
	}

	if strings.HasPrefix(input, "disconnect ") {
		indexString := strings.TrimPrefix(input, "disconnect ")
		index, err := strconv.Atoi(indexString)
		if err != nil {
			return
		}
		for _, r := range s.connectedPlayersList {
			if r.id == index {
				r.conn.Close()
				return
			}
		}

		fmt.Println("Index not found")
	}
}

func requestClientList(c client, s *Server) {
	sendClientList(c.id, s)
}

func disconnectClient(c client, s *Server) {
	c.conn.Close()
}
