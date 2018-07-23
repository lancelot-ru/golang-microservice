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
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("X-Custom-Header", "fromServer")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&m)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	return m
}

func handleJSON(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var m Message
	err := decoder.Decode(&m)
	if err != nil {
		log.Println(err)
	}
	if m.System == "A" {
		fl := findAId(m.OperationID, m.Action)
		if fl == true {
			log.Println("Already in the table!")
			m.System = "SRV"
			m.OperationID = 0
			m.Action = "Already in the table!"
		} else {
			insertNewAction(m.OperationID, m.Action)
			log.Println("Sending to B...")
			m = answerSystemsHTTP("http://localhost:8083/b", m)
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
	r := http.NewServeMux()
	r.HandleFunc("/", handleJSON)

	log.Fatal(http.ListenAndServe(":8080", r))
}
