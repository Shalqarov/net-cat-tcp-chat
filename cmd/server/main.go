package main

import (
	"flag"
	"log"
	"os"

	"github.com/Shalqarov/net-cat/internal"
)

var (
	listen = flag.Bool("l", false, "Listen")
	host   = flag.String("h", "localhost", "Host")
	port   = flag.Int("p", 0, "Port")
)

func main() {
	flag.Parse()
	if !*listen {
		log.Println("Listen flag is not true")
		return
	}
	file, err := os.OpenFile("log/log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	defer file.Close()
	internal.StartServer(*host, *port)
}
