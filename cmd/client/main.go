package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Shalqarov/net-cat/pkg"
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatalln("Hostname and port required")
	}
	serverHost := flag.Arg(0)
	serverPort := flag.Arg(1)
	pkg.StartClient(fmt.Sprintf("%s:%s", serverHost, serverPort))
}
