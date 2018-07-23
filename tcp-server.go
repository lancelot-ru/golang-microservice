// tcp-server
package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"

	"golang-microservice/proto"

	"github.com/golang/protobuf/proto"
)

func answerSystemsTCP(address string, m Message) Message {
	msg1 := &msg.Message{
		System:      m.System,
		OperationId: int32(m.OperationID),
		Action:      m.Action,
	}

	conn, _ := net.Dial("tcp", address)
	data, err := proto.Marshal(msg1)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	conn.Write(data)

	// listen for reply
	var buf bytes.Buffer
	log.Println("Receiving data…")
	_, err = io.Copy(&buf, conn)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	pdata := new(msg.Message)
	err = proto.Unmarshal(buf.Bytes(), pdata)
	if err != nil {
		log.Println(err)
	}
	m.System = pdata.GetSystem()
	m.OperationID = int(pdata.GetOperationId())
	m.Action = pdata.GetAction()
	return m
}

func handleReceivedData(data *msg.Message) {
	log.Println("Handling data...")
	var m Message
	m.System = data.GetSystem()
	m.OperationID = int(data.GetOperationId())
	m.Action = data.GetAction()
	log.Println(m.System, m.OperationID, m.Action)
	if m.System == "A" {
		fl := findAId(m.OperationID, m.Action)
		if fl == true {
			log.Println("Already in the table!")
		} else {
			insertNewAction(m.OperationID, m.Action)
			log.Println("Sending to B...")
		}

		showAllRows()
	}
	if m.System == "B" {
		fl := updateBId(m.OperationID, m.Action, false)
		if fl == true {
			log.Println("Success!")
		} else {
			log.Println("Error!")
		}
		log.Println("Sending to A...")
		showAllRows()
	}
}

func handleProtoClient(conn net.Conn) {
	log.Println("Connected!")
	defer conn.Close()
	var buf bytes.Buffer
	log.Println("Receiving data…")
	_, err := io.Copy(&buf, conn)
	if err != nil {
		log.Println(err)
	}
	pdata := new(msg.Message)
	err = proto.Unmarshal(buf.Bytes(), pdata)
	if err != nil {
		log.Println(err)
	}

	handleReceivedData(pdata)
}

func serveTCP() {
	log.Println("Staring Server..")

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			go handleProtoClient(conn)
		} else {
			continue
		}
	}
}
