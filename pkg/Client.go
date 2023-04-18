package pkg

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func (server *Server) Client(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		}
		if server.Counter == 10 {
			conn.Write([]byte("Chat is full. Try next time"))
			conn.Close()
			continue
		}
		server.Counter++

		go server.HandleClient(conn)
	}
}

func (server *Server) HandleClient(conn net.Conn) {
	conn.Write([]byte("Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'"))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig, "erger")
		}
	}()
	for {
		fmt.Println("sleeping...")
		time.Sleep(2 * time.Second)
	}
	server.Counter--
	conn.Close()
	fmt.Println("someone leaved chat")
}
