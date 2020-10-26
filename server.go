package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error, clientid int, msgs chan Message) {
	// TODO: all
	// Deal with an error event.
	fmt.Println(err.Error())
	// I think this means a user has disconnected
	if err.Error() == "EOF" {
		msgs <- Message{-2, strconv.Itoa(clientid)}
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, _ := ln.Accept()
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	msgs <- Message{-1, strconv.Itoa(clientid)}
	reader := bufio.NewReader(client)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			go handleError(err, clientid, msgs)
			break
		}
		fmt.Printf(msg)
		msgTidied := Message{clientid, msg}
		msgs <- msgTidied
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, _ := net.Listen("tcp", *portPtr)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	clientId := 0
	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client
			clients[clientId] = conn
			go handleClient(conn, clientId, msgs)
			clientId += 1
		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			for id, conn := range clients {
				if msg.sender == -1 {
					fmt.Fprintf(conn, "User %s has joined the room\n", msg.message)
				} else if msg.sender == -2 {
					fmt.Fprintf(conn, "User %s has disconected\n", msg.message)
					userIdInteger, _ := strconv.Atoi(msg.message)
					delete(clients, userIdInteger)
				} else if msg.sender != id {
					fmt.Fprintf(conn, "User %d: %s", msg.sender, msg.message)
				}
			}
		}
	}
}
