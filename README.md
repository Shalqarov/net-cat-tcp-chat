## Net-cat

## How to start
### Server:
```console
[USAGE]: go run ./cmd/server/ *port*
$ go run ./cmd/server/ 8989
Listening on the port :8989
$ go run ./cmd/server/ 2525
Listening on the port :2525
$ go run ./cmd/server/ 2525 localhost
[USAGE]: ./TCPChat $port
$
```

### Client:
```
$ go run ./cmd/client *host* *port*
for example
$ go run ./cmd/client localhost 8989
```

## Objectives

This project consists on recreating the **NetCat in a Server-Client Architecture** that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

- NetCat, `nc` system command, is a command-line utility that reads and writes data across network connections using TCP or UDP. It is used for anything involving TCP, UDP, or UNIX-domain sockets, it is able to open TCP connections, send UDP packages, listen on arbitrary TCP and UDP ports and many more.

- To see more information about NetCat inspect the manual `man nc`.

This project must work in a similar way that the original  NetCat works, in other words, you must create a group chat. The project must have the following features :

- TCP connection between server and multiple clients (relation of 1 to many).
- A name requirement to the client.
- Control connections quantity.
- Clients must be able to send messages to the chat.
- Do not broadcast EMPTY messages from a client.
- Messages sent, must be identified by the time that was sent and the user name of who sent the message, example : `[2020-01-20 15:48:41][client.name]:[client.message]`
- If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.
- If a Client connects to the server, the rest of the Clients must be informed by the server that the Client joined the group.
- If a Client exits the chat, the rest of the Clients must be informed by the server that the Client left.
- All Clients must receive the messages sent by other Clients.
- If a Client leaves the chat, the rest of the Clients must not disconnect.
- If there is no port specified, then set as default the port 8989. Otherwise, program must respond with usage message: `[USAGE]: ./TCPChat $port`

### Instructions

- Project must be written in **Go**
- Start TCP server, listen and accept connections
- Project must have Go-routines
- Project must have channels or Mutexes
- Maximum 10 connections
- The code must respect the [**good practices**](../good-practices/README.md)
- It is recommended to have **test files** for [unit testing](https://go.dev/doc/tutorial/add-a-test) both the server connection and the client.

### Allowed Packages

- io
- log
- os
- fmt
- net
- sync
- time
- bufio
- errors
- strings
- reflect

## This project will help learn about :

- Manipulation of structures.
- [Net-Cat](https://www.commandlinux.com/man-page/man1/nc.1.html)
- TCP/UDP
  - TCP/UDP connection
  - TCP/UDP socket
- [Go concurrency](https://golang.org/doc/#go_concurrency_patterns)
  - [Channels](https://tour.golang.org/concurrency/2)
  - [Goroutines](https://tour.golang.org/concurrency/1)
- Mutexes
- IP and [ports](https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers)
