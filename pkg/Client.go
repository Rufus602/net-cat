package pkg

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func (server *Server) Client(l net.Listener) {
	var maxNumCon int = 10
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		}
		if server.Counter > maxNumCon {
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
	var name string
	for {
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		temp, err := readFromClient(conn)
		if err != nil {
			return
		}
		if _, ok := server.Clients[name]; ok {
			fmt.Fprintln(conn, "Username already in use")
			continue
		}
		name = temp
		break
	}
	loadHistory(conn)
	server.Clients[name] = conn
	server.Join <- name
	for {
		fmt.Fprintf(conn, "[%s][%s]: ", time.Now().Format("01-02-2006 15:04:05"), name)

		msg, err := readFromClient(conn)
		if err != nil {
			server.Mu.Lock()
			server.Left <- name
			server.Counter--
			server.Mu.Unlock()
			break
		}
		m := Message{
			Sender: name,
			Text:   msg,
		}
		server.Message <- m
	}
	server.Mu.Lock()

	delete(server.Clients, name)
	server.Mu.Unlock()
}

func readFromClient(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return "", err
	}
	str := strings.TrimSpace(string(buf[:n]))
	return str, nil
}

func loadHistory(conn net.Conn) {
	file, err := os.Open("History.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintln(conn, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		return
	}
}
