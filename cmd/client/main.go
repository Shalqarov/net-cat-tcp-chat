package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatalln("Hostname and port required")
	}
	serverHost := flag.Arg(0)
	serverPort := flag.Arg(1)
	startClient(fmt.Sprintf("%s:%s", serverHost, serverPort))
}

func startClient(addr string) {
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
