package pkg

type Server struct {
	Counter int
	Clients []string
	Message chan string
}
