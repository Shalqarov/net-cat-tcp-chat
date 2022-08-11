package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	listen = flag.Bool("l", false, "Listen")
	host   = flag.String("h", "localhost", "Host")
	port   = flag.Int("p", 0, "Port")
)

const (
	greetings = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "

	timeFormat = "2006-01-02 15:04:05"
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
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not listen: %v", err)
	}
	defer listener.Close()

	log.Printf("Listening for connections on %s", listener.Addr().String())
	connections := Connect{
		users: &sync.Map{},
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
			conn.Close()
			continue
		}
		connections.register(conn)
		go connections.handleConnection(conn)
	}
}

func (c *Connect) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		c.users.Delete(conn)
	}()

	for {
		userInput, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}
		if userInput == "\n" {
			fmt.Fprint(conn, c.message(time.Now(), conn))
			continue
		}
		c.users.Range(func(key, value interface{}) bool {
			if _, ok := value.(string); ok && key.(net.Conn) != conn {
				c.sendMessage(key.(net.Conn), userInput)
			}
			return true
		})
		fmt.Fprint(conn, c.message(time.Now(), conn))
	}
}

func (c *Connect) message(now time.Time, conn net.Conn) string {
	username, _ := c.users.Load(conn)
	return "[" + now.Format(timeFormat) + "]" + "[" + username.(string) + "]:"
}

func (c *Connect) sendMessage(conn net.Conn, input string) {
	fmt.Fprintln(conn)
	fmt.Fprint(conn, c.message(time.Now(), conn))
	fmt.Fprint(conn, input)
	fmt.Fprint(conn, c.message(time.Now(), conn))
}

func (c *Connect) register(conn net.Conn) {
	username := ""
	fmt.Fprint(conn, greetings)
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.users.Store(conn, strings.TrimRight(username, "\r\n"))
	fmt.Fprint(conn, c.message(time.Now(), conn))
}
