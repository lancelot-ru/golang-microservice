// tcp-server
package main

import (
	"log"
	"net"
	"os"

	"golang-microservice/proto"

	"github.com/golang/protobuf/proto"
)

func answerSystemsTCP(address string, m Message) Message {
	// send
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

	checkError(err, "listening")
	log.Println("Decoding Protobuf message...")
	pdata := new(msg.Message)
	err = proto.Unmarshal(data2[0:n], pdata)
	checkError(err, "unmarshalling")

	m.System = pdata.GetSystem()
	m.OperationID = int(pdata.GetOperationId())
	m.Action = pdata.GetAction()

	log.Println(m.System, m.OperationID, m.Action)
	return m
}

func handleReceivedData(data *msg.Message) Message {
	log.Println("Handling data...")
	var m Message
	m.System = data.GetSystem()
	m.OperationID = int(data.GetOperationId())
	m.Action = data.GetAction()

	log.Println(m.System, m.OperationID, m.Action)

	if m.System == "A" {
		fl := findAId(m.OperationID, m.Action)
		if fl == 1 {
			log.Println("Already in the table!")
			m.System = "SRV"
			m.OperationID = 0
			m.Action = "Already in the table!"
		} else if fl == 0 {
			insertNewAction(m.OperationID, m.Action)
			log.Println("Sending to B...")
			m = answerSystemsTCP("localhost:8083", m)
		} else {
			log.Println("Error!")
			m.System = "SRV"
			m.OperationID = 0
			m.Action = "Error!"
		}

		showAllRows()
	}
	if m.System == "B" {
		fl := updateBId(m.OperationID, m.Action, false)
		if fl == true {
			log.Println("Success!")
			m.System = "SRV"
			m.OperationID = 0
			m.Action = "Success!"
		} else {
			log.Println("Error!")
			m.System = "SRV"
			m.OperationID = 0
			m.Action = "Error!"
		}

		showAllRows()

		log.Println("Sending to A...")
	}
	return m
}

func handleProtoMessage(conn net.Conn) {
	log.Println("Connected!")
	defer conn.Close()
	data := make([]byte, 4096)
	n, err := conn.Read(data)
	checkError(err, "reading")
	log.Println("Decoding Protobuf message...")
	pdata := new(msg.Message)
	err = proto.Unmarshal(data[0:n], pdata)
	checkError(err, "unmarshalling")

	m := handleReceivedData(pdata)

	msg2 := &msg.Message{
		System:      m.System,
		OperationId: int32(m.OperationID),
		Action:      m.Action,
	}
	log.Println("Encoding it back...")
	data2, err := proto.Marshal(msg2)
	checkError(err, "marshalling")
	conn.Write(data2)
}

func serveTCP() {
	log.Println("Staring TCP Server..")

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			go handleProtoMessage(conn)
		} else {
			continue
		}
	}
}
