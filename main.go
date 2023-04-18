package main

import (
	"fmt"
	"net"
	. "net-cat/pkg"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
	} else if len(args) == 0 {
		args = append(args, "8989")
	}
	l, err := net.Listen("tcp", ":"+args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("Listening on localhost:" + args[0])
	defer l.Close()
	var server Server = Server{0, nil, nil}
	server.Client(l)
}
