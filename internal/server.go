package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	sync.Mutex
	users     *sync.Map
	connCount uint8
	history   string
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.logout(conn)
	s.Lock()
	s.connCount++
	s.Unlock()
	err := s.register(conn)
	if err != nil {
		fmt.Fprintf(conn, "Registration failed: %v\n", err.Error())
		return
	}
	for {
		userInput, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		if userInput == "\n" {
			fmt.Fprint(conn, s.message(conn))
			continue
		}
		if userInput == "--exit\n" {
			break
		}
		s.sendMessage(conn, userInput)
		fmt.Fprint(conn, s.message(conn))
	}
	log.Printf("Closing conection with %s", conn.RemoteAddr())
}

func (s *Server) sendMessage(conn net.Conn, input string) {
	s.Lock()
	s.history += s.message(conn) + input
	s.Unlock()
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
