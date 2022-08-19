package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func readRoutine(ctx context.Context, cancel context.CancelFunc, conn net.Conn) {
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				cancel()
				break OUTER
			}
			fmt.Print(string(buff[0:n]))
		}
	}
	log.Printf("Exiting...")
}

func writeRoutine(ctx context.Context, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				break OUTER
			}
			str := scanner.Text()
			conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}
	}
}

func main() {
	flag.Parse()
	dialer := &net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := fmt.Sprintf("%s:%s", host, port)

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		readRoutine(ctx, cancel, conn)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		writeRoutine(ctx, conn)
		wg.Done()
	}()

	wg.Wait()
	conn.Close()
}
