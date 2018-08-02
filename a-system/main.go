// a-system-main
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"

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

func sendMessage(url string, m Message) {
	b, err := json.Marshal(m)
	checkError(err, "marshalling")
	log.Println(string(b))

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("X-Custom-Header", "fromServer")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err, "client")
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&m)
	checkError(err, "decoding JSON")
	if m.System == "SRV" {
		log.Println(m.Action)
	} else {
		log.Println("Error getting response from the server")
	}
	defer resp.Body.Close()
}

func sendMessageTCP(address string, m Message) {
	//send
	msg1 := &msg.Message{
		System:      m.System,
		OperationId: int32(m.OperationID),
		Action:      m.Action,
	}

	conn, _ := net.Dial("tcp", address)
	defer conn.Close()
	data1, err := proto.Marshal(msg1)
	checkError(err, "marshalling")
	conn.Write(data1)

	// listen for reply
	data2 := make([]byte, 4096)
	n, err := conn.Read(data2)
	checkError(err, "reading")
	log.Println("Decoding Protobuf message...")

	pdata := new(msg.Message)
	err = proto.Unmarshal(data2[0:n], pdata)
	checkError(err, "unmarshalling")

	m.System = pdata.GetSystem()
	m.OperationID = int(pdata.GetOperationId())
	m.Action = pdata.GetAction()

	log.Println(m.System, m.OperationID, m.Action)
}

func main() {
	log.Println("A system started")
	s := "tcp" //or "http"
	if s == "http" {
		initial := Message{"A", 11, "land"}
		sendMessage("http://localhost:8080", initial)
	}
	if s == "tcp" {
		initial := Message{"A", 17, "gaze"}
		sendMessageTCP("localhost:8082", initial)
	}
}
