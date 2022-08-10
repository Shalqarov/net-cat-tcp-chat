package pkg

import (
	"fmt"
	"log"
	"net"
	"time"
)

func StartClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("could not connect to the server: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connected...")
	for {
		go handleConnection(conn)
		var source string
		_, err := fmt.Scanln(&source)
		if err != nil {
			fmt.Println("invalid input")
			continue
		}
		if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Println()
	}
}

func handleConnection(conn net.Conn) {
	for {
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			break
		}
		fmt.Print(string(buff[0:n]))
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 700))
	}
}
