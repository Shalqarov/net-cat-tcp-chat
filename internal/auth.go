package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const WELCOME_MSG = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "

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

func (s *Server) logout(conn net.Conn) {
	s.logoutMessage(conn)
	conn.Close()
	s.Lock()
	s.connCount--
	s.Unlock()
	s.users.Delete(conn)
}

func (s *Server) welcomeMessage(conn net.Conn) {
	username, _ := s.users.Load(conn)
	msg := fmt.Sprintf("%s has joined our chat...\n", username.(string))
	fmt.Printf("%s has joined to the server\n", username.(string))
	log.Printf("%s has joined to the server\n", username.(string))
	s.notification(conn, msg)
}

func (s *Server) logoutMessage(conn net.Conn) {
	username, _ := s.users.Load(conn)
	msg := fmt.Sprintf("%s has left our chat...\n", username.(string))
	fmt.Printf("%s has left server\n", username.(string))
	log.Printf("%s has left server\n", username.(string))
	s.notification(conn, msg)
}
