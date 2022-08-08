package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	listen = flag.Bool("l", false, "Listen")
	host   = flag.String("h", "localhost", "Host")
	port   = flag.Int("p", 0, "Port")
)

func main() {
	flag.Parse()
	if *listen {
		startServer()
		return
	}
}

func startServer() {
	addr := fmt.Sprintf("%s:%d", *host, *port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not listen: %v", err)
	}

	log.Printf("Listening for connections on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
			continue
		}
		go processClient(conn)
	}
}

func processClient(conn net.Conn) {
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		log.Println(err)
	}
	conn.Close()
}
