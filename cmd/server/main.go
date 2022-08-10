package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	listen = flag.Bool("l", false, "Listen")
	host   = flag.String("h", "localhost", "Host")
	port   = flag.Int("p", 0, "Port")
)

const (
	greetings = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "
)

type Connect struct {
	sync.Mutex
	users *sync.Map
}

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
	connections := Connect{
		users: &sync.Map{},
	}
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

		go connections.handleConnection(conn)
	}
}

func (c *Connect) handleConnection(conn net.Conn) {
	username := ""
	fmt.Fprint(conn, greetings)
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.users.Store(username, conn)

	defer func() {
		conn.Close()
		c.users.Delete(username)
	}()

	for {
		userInput, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println(userInput)
		c.users.Range(func(key, value interface{}) bool {
			if user, ok := value.(net.Conn); ok && user != conn {
				if _, err := conn.Write([]byte(userInput)); err != nil {
					log.Println("error on writing to connection", err.Error())
				}
			}
			return true
		})
	}
}
