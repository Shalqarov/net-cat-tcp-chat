package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Shalqarov/net-cat/internal"
)

func main() {
	if _, err := os.Stat("log"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("log", os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}
	file, err := os.OpenFile("log/log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	defer file.Close()

	args := os.Args[1:]
	port, err := internal.PortParse(args)
	if err != nil {
		log.Println(err)
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	internal.StartServer(port)
}
