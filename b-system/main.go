// b-system-main
package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"golang-microservice/proto"

	"github.com/golang/protobuf/proto"
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

func handleHTTP(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var m Message
	err := decoder.Decode(&m)
	checkError(err, "decoding JSON")

	var newMessage Message
	newMessage.System = "B"
	newMessage.OperationID = m.OperationID + 1
	newMessage.Action = m.Action

	log.Println(newMessage.System, newMessage.OperationID, newMessage.Action)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(newMessage)
}

func handleTCP(conn net.Conn) {
	log.Println("Connected!")
	defer conn.Close()
	data1 := make([]byte, 4096)
	n, err := conn.Read(data1)
	checkError(err, "reading Protobuf")
	log.Println("Decoding Protobuf message")
	pdata := new(msg.Message)
	err = proto.Unmarshal(data1[0:n], pdata)
	checkError(err, "unmarshalling")

	var m Message
	m.System = pdata.GetSystem()
	m.OperationID = int(pdata.GetOperationId())
	m.Action = pdata.GetAction()

	log.Println(m.System, m.OperationID, m.Action)

	msg2 := &msg.Message{
		System:      "B",
		OperationId: int32(m.OperationID) + 1,
		Action:      m.Action,
	}
	log.Println("Encoding it back...")
	data2, err := proto.Marshal(msg2)
	checkError(err, "marshalling")
	conn.Write(data2)
}

func main() {
	log.Println("B system started")
	s := "tcp" //or "http"
	if s == "tcp" {
		listener, err := net.Listen("tcp", ":8083")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		for {
			if conn, err := listener.Accept(); err == nil {
				go handleTCP(conn)
			} else {
				continue
			}
		}
	}
	if s == "http" {
		r := http.NewServeMux()
		r.HandleFunc("/b", handleHTTP)

		log.Fatal(http.ListenAndServe(":8083", r))
	}
}
