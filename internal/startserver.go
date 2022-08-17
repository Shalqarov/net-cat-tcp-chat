package internal

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	TIME_FORMAT     = "2006-01-02 15:04:05"
	MAX_CONNECTIONS = uint8(10)
)

func StartServer(port int) {
	addr := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not listen: %v", err)
	}
	defer listener.Close()

	log.Printf("Listening on the port :%d\n", port)
	fmt.Printf("Listening on the port :%d\n", port)
	server := Server{
		users: &sync.Map{},
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
			conn.Close()
			continue
		}
		server.Lock()
		if server.connCount >= MAX_CONNECTIONS {
			fmt.Fprintf(conn, "Connection failed, chat room is full.\nMax connection number is %d\n", MAX_CONNECTIONS)
			conn.Close()
		}
		server.Unlock()
		go server.handleConnection(conn)
	}
}
