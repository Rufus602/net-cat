package pkg

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Server struct {
	Counter int
	Clients map[string]net.Conn
	Message chan Message
	Join    chan string
	Left    chan string
	Mu      *sync.Mutex
}

type Message struct {
	Sender string
	Text   string
}

func (server *Server) Spreader() {
	file, err := os.OpenFile("History.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	for {
		current_time := time.Now().Format("01-02-2006 15:04:05")
		select {
		case j := <-server.Join:
			server.Mu.Lock()
			for name, conn := range server.Clients {
				if name == j {
					continue
				}
				// file.WriteString(j + " has joined the chat!\n")
				fmt.Fprintln(conn, "\n"+j+" has joined the chat!")
				fmt.Fprintf(conn, "[%s][%s]: ", current_time, name)
			}
			server.Mu.Unlock()
		case l := <-server.Left:
			server.Mu.Lock()
			for name, conn := range server.Clients {
				if name == l {
					continue
				}
				// file.WriteString(l + " has left the chat!\n")
				fmt.Fprintln(conn, "\n"+l+" has left the chat!")
				fmt.Fprintf(conn, "[%s][%s]: ", current_time, name)
			}
			server.Mu.Unlock()
		case m := <-server.Message:
			server.Mu.Lock()
			if m.Text == "" {
				continue
			}
			file.WriteString(fmt.Sprintf("[%s][%s]: %s\n", current_time, m.Sender, m.Text))
			for name, conn := range server.Clients {
				if name == m.Sender {
					continue
				}
				// file.WriteString(fmt.Sprintf("\n[%s][%s]: %s\n", current_time, m.sender, m.text))
				fmt.Fprintf(conn, "\n[%s][%s]: %s\n", current_time, m.Sender, m.Text)
				fmt.Fprintf(conn, "[%s][%s]: ", current_time, name)
			}
			server.Mu.Unlock()
		}
	}
}
