package main

import (
	"log"
)

type Message struct {
	System      string
	OperationID int
	Action      string
}

func checkError(err error, place string) {
	if err != nil {
		log.Println("Error", err, "at", place)
	}
}

func main() {
	go func() { serveHTTP() }()
	go func() { serveTCP() }()
	select {}
}
