package main

import (
	"flag"
	"log"

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
	internal.StartServer(*host, *port)
}
