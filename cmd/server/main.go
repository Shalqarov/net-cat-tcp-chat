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
	startServer(*host, *port)
}

func startServer(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
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
		server.Lock()
		if server.connCount >= MAX_CONNECTIONS {
			fmt.Fprintf(conn, "Connection failed, chat room is full.\nMax connection number is %d\n", MAX_CONNECTIONS)
			conn.Close()
		}
		server.Unlock()
		go server.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.logout(conn)
	err := s.register(conn)
	if err != nil {
		fmt.Fprintf(conn, "Registration failed: %v\n", err.Error())
		return
	}
	for {
		userInput, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}
		if userInput == "\n" {
			fmt.Fprint(conn, s.message(conn))
			continue
		}
		s.sendMessage(conn, userInput)
		fmt.Fprint(conn, s.message(conn))
	}
}

func (s *Server) logout(conn net.Conn) {
	s.logoutMessage(conn)
	conn.Close()
	s.Lock()
	s.connCount--
	s.Unlock()
	s.users.Delete(conn)
}

func (s *Server) message(conn net.Conn) string {
	username, _ := s.users.Load(conn)
	return "[" + time.Now().Format(TIME_FORMAT) + "]" + "[" + username.(string) + "]:"
}

func (s *Server) welcomeMessage(conn net.Conn) {
	username, _ := s.users.Load(conn)
	msg := fmt.Sprintf("%s has joined our chat...\n", username.(string))
	s.notification(conn, msg)
}

func (s *Server) notification(conn net.Conn, msg string) {
	s.users.Range(func(key, value interface{}) bool {
		if _, ok := value.(string); ok && key.(net.Conn) != conn {
			fmt.Fprintln(key.(net.Conn))
			fmt.Fprint(key.(net.Conn), msg)
			fmt.Fprint(key.(net.Conn), s.message(key.(net.Conn)))
		}
		return true
	})
}

func (s *Server) logoutMessage(conn net.Conn) {
	username, _ := s.users.Load(conn)
	msg := fmt.Sprintf("%s has left our chat...\n", username.(string))
	s.notification(conn, msg)
}

func (s *Server) sendMessage(conn net.Conn, input string) {
	s.users.Range(func(key, value interface{}) bool {
		if _, ok := value.(string); ok && key.(net.Conn) != conn {
			fmt.Fprintln(key.(net.Conn))
			fmt.Fprint(key.(net.Conn), s.message(conn))
			fmt.Fprint(key.(net.Conn), input)
			fmt.Fprint(key.(net.Conn), s.message(key.(net.Conn)))
		}
		return true
	})
}

func (s *Server) register(conn net.Conn) error {
	fmt.Fprint(conn, WELCOME_MSG)
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}
	s.Lock()
	s.connCount++
	s.Unlock()
	s.users.Store(conn, strings.TrimRight(username, "\r\n"))
	fmt.Fprint(conn, s.message(conn))
	s.welcomeMessage(conn)
	return nil
}
