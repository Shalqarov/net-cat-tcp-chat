package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Shalqarov/net-cat/internal"
)

func PortParse(args []string) (int, error) {
	if len(args) > 1 {
		return 0, errors.New("invalid number of args")
	} else if len(args) == 0 {
		return 8989, nil
	}
	return strconv.Atoi(args[0])
}

func main() {
	file, err := os.OpenFile("log/log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	defer file.Close()

	args := os.Args[1:]
	port, err := PortParse(args)
	if err != nil {
		log.Println(err)
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	internal.StartServer(port)
}
