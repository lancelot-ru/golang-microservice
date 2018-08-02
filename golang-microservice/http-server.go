// http-server
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func answerSystemsHTTP(url string, m Message) Message {
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

	defer resp.Body.Close()
	return m
}

func handleJSON(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var m Message
	err := decoder.Decode(&m)
	checkError(err, "decoding JSON")

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
			m = answerSystemsHTTP("http://localhost:8083/b", m)
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
		log.Println("Sending to A...")
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
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(m)
}

func serveHTTP() {
	log.Println("Staring HTTP Server..")

	r := http.NewServeMux()
	r.HandleFunc("/", handleJSON)

	log.Fatal(http.ListenAndServe(":8080", r))
}
