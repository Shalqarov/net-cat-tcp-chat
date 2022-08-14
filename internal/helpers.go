package internal

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

func PortParse(args []string) (int, error) {
	if len(args) > 1 {
		return 0, errors.New("invalid number of args")
	} else if len(args) == 0 {
		return 8989, nil
	}
	return strconv.Atoi(args[0])
}

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
