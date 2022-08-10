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

const (
	greetings      = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "
	timeFormat     = "2006-01-02 15:04:05"
	maxConnections = 10
	usage          = "[USAGE]: ./TCPChat $port"
)

func main() {
	flag.Parse()
	if !*listen {
		log.Println("Listen flag is not true")
		return
	}
	startServer()
}

func startServer() {
	addr := fmt.Sprintf("%s:%d", *host, *port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not listen: %v", err)
	}
	defer listener.Close()

	log.Printf("Listening for connections on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
			conn.Close()
			continue
		}
		fmt.Fprint(conn, greetings)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		log.Println(err)
	}
}
