package pkg

import (
	"io"
	"log"
	"net"
	"os"
)

func StartClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("could not connect to the server: %v", err)
	}
	defer conn.Close()
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		log.Fatalf("Connection error: %s", err)
	}
}
