package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
