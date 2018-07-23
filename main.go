package main

type Message struct {
	System      string
	OperationID int
	Action      string
}

func main() {
	go func() { serveHTTP() }()
	go func() { serveTCP() }()
	select {}
}
