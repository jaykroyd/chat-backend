package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

var port string = ":8080"
var activeServerList []*Server

type client struct {
	username string
	id       int
	conn     net.Conn
	addr     net.Addr
}

// Server struct contains a channel that's used to signal shutdown and a wait group to wait until all the server's goroutines are actually done
type Server struct {
	maxPlayers           int
	connectedPlayersList []client
	listener             net.Listener
	quit                 chan interface{}
	wg                   sync.WaitGroup
}

func main() {
	// Initialize active server list
	activeServerList = make([]*Server, 0)

	// Attempt to start new server
	NewServer(port, 4)
}

// NewServer creates and returns a new server
func NewServer(addr string, maxPlayers int) *Server {
	println("Starting server on port: " + port)

	s := &Server{
		maxPlayers:           maxPlayers,
		connectedPlayersList: make([]client, maxPlayers),
		quit:                 make(chan interface{}),
	}

	activeServerList = append(activeServerList, s)
	fmt.Println(s.connectedPlayersList)

	// Start listening for incoming tcp connections
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("fatal error")
		log.Fatal(err)
	}

	s.listener = ln
	s.wg.Add(1)

	go waitForInput(s)
	s.serve()

	return s
}

func (s *Server) serve() {
	// Start loop to listen for new connections
	for {
		// Accept incoming connection
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				fmt.Println("Closing listner")
				return
			default:
				log.Println("Error accepting new connection.", err)
			}
		} else {
			// Check if server is full
			slot := assignServerSlot(s)
			if slot == -1 {
				fmt.Println("Server is full!")
				conn.Close()
				continue
			}
			// Increment number of connected clients
			c := client{
				"player " + strconv.Itoa(1+slot),
				slot,
				conn,
				conn.RemoteAddr(),
			}
			if err := addConnectedClient(c, s); err != nil {
				conn.Close()
				panic(err)
			}

			// Add +1 connection to wait group
			s.wg.Add(1)

			go func() {
				// Start listening for input
				s.handleConnection(c)
				s.wg.Done()
			}()

			// Send welcome packet
			sendWelcome(c.id, s, "Welcome to the server!")
		}

	}
}

func (s *Server) handleConnection(c client) {
	// Defer the closure of the connection, to ensure resources are released once the loop is finished
	defer c.conn.Close()

	// Start loop to listen for new input from existing connection
	for {
		err := readDataFromClient(c.conn, s)
		if err != nil {
			fmt.Println(err)
			removeConnectedClient(c.conn, s)
			return
		}
	}
}

func waitForInput(s *Server) {
	// Create scanner to read input in real time
	scanner := bufio.NewScanner(os.Stdin)

	// Initial message
	fmt.Println("Please type in your input:")

	// Start loop to listen for input from client
	for scanner.Scan() {
		input := scanner.Text()
		fmt.Println("Input received from server:", input)

		checkForInput(input, s)
	}
}

// Stop orderly shuts down the server
func (s *Server) Stop() {
	fmt.Println("Stopping server...")
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}
