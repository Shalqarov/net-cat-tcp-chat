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
	WELCOME_MSG     = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "
	TIME_FORMAT     = "2006-01-02 15:04:05"
	MAX_CONNECTIONS = uint8(10)
)

type Server struct {
	sync.Mutex
	users     *sync.Map
	connCount uint8
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
		if server.connCount >= MAX_CONNECTIONS {
			fmt.Fprintf(conn, "Connection failed, chat room is full.\nMax connection number is %d\n", MAX_CONNECTIONS)
			conn.Close()
		}
		server.register(conn)
		go server.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.Lock()
		s.connCount--
		s.Unlock()
		s.users.Delete(conn)
	}()

	for {
		userInput, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}
		if userInput == "\n" {
			fmt.Fprint(conn, s.message(conn))
			continue
		}
		s.users.Range(func(key, value interface{}) bool {
			if _, ok := value.(string); ok && key.(net.Conn) != conn {
				s.sendMessage(key.(net.Conn), userInput)
			}
			return true
		})
		fmt.Fprint(conn, s.message(conn))
	}
}

func (s *Server) message(conn net.Conn) string {
	username, _ := s.users.Load(conn)
	return "[" + time.Now().Format(TIME_FORMAT) + "]" + "[" + username.(string) + "]:"
}

func (s *Server) sendMessage(conn net.Conn, input string) {
	fmt.Fprintln(conn)
	fmt.Fprint(conn, s.message(conn))
	fmt.Fprint(conn, input)
	fmt.Fprint(conn, s.message(conn))
}

func (s *Server) register(conn net.Conn) {
	fmt.Fprint(conn, WELCOME_MSG)
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err.Error())
		return
	}
	s.Lock()
	s.connCount++
	s.Unlock()
	s.users.Store(conn, strings.TrimRight(username, "\r\n"))
	fmt.Fprint(conn, s.message(conn))
}
