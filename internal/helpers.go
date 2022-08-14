package internal

import (
	"fmt"
	"net"
	"time"
)

func (s *Server) message(conn net.Conn) string {
	username, _ := s.users.Load(conn)
	return "[" + time.Now().Format(TIME_FORMAT) + "]" + "[" + username.(string) + "]:"
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
