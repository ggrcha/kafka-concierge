package transaction

import (
	"encoding/json"
	"io/ioutil"
	debuggin "kernel-concierge/Debuggin"
	kafka "kernel-concierge/Kafka"
	pending "kernel-concierge/Pending"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

// ResponseMessage ...
type ResponseMessage struct {
	Msg string `json:"message"`
}

// RequestData ...
type RequestData struct {
	IDRequest    string `json:"idRequest"`
	JaggerParams string `json:"jagerParams"`
	// Accounts     struct {
	// 	AccountsOperations []struct {
	// 		ID    string `json:"id"`
	// 		Value int    `json:"value"`
	// 	} `json:"accountsOperations"`
	// }
	AccountsOperations string
}

// const letterBytes = "ab"

var rm ResponseMessage

// Manager manages request to kafka
func Manager(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	rm = ResponseMessage{}

	// rand.Seed(time.Now().UnixNano())

	// rest boilerplate
	rd := RequestData{}
	req, _ := ioutil.ReadAll(r.Body)

	rd.AccountsOperations = strings.TrimSpace(string(req))
	rd.IDRequest = uuid.Must(uuid.NewV4()).String()

	// if err != nil {
	// 	log.Println(debuggin.Tracer(), "invalid account operations", err)
	// 	rm.Msg = "Invalid account operations format"
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	resp, _ := json.Marshal(rm)
	// 	w.Write(resp)
	// 	return
	// }

	// creates pending request to add do the stream pending map
	pr := pending.Request{}
	// pr.RequestID = string(letterBytes[rand.Intn(len(letterBytes))])
	pr.RequestID = rd.IDRequest
	pr.ResponseChan = make(chan string)
	pr.RequestData = rd.AccountsOperations
	// defers closing the channel and any other resource opened
	defer closeResources(pr)

	// adds newly created request to control map
	kafka.NewRequest <- pr

	kafka.ProduceRequest(pr)

	// blocks execution until response or timeout
	select {
	case <-time.After(3 * time.Second):
		log.Println(debuggin.Tracer(), "timeout received")
		// Notifies timeout
		kafka.ToChan <- pr
	case <-pr.ResponseChan:
		// returns response to client
		log.Println(debuggin.Tracer(), "received response from kafka consumer")
		rm.Msg = "success"
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(rm)
		w.Write(resp)
		log.Println(debuggin.Tracer(), "postei a mensagem no kafka")
	}
}

func closeResources(pr pending.Request) {
	close(pr.ResponseChan)
}
