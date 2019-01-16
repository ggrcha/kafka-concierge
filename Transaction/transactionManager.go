package transaction

import (
	"encoding/json"
	"io/ioutil"
	debuggin "kernel-concierge/Debuggin"
	kafka "kernel-concierge/Kafka"
	pending "kernel-concierge/Pending"
	"log"
	"net/http"
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
	Accounts     struct {
		AccountsOperations []struct {
			ID    string `json:"id"`
			Value int    `json:"value"`
		} `json:"accountsOperations"`
	} `json:"accounts"`
}

var rm ResponseMessage

// Manager manages request to kafka
func Manager(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	rm = ResponseMessage{}

	// rest boilerplate
	rd := RequestData{}
	req, _ := ioutil.ReadAll(r.Body)

	_ = json.Unmarshal(req, &rd.Accounts)
	rd.IDRequest = uuid.Must(uuid.NewV4()).String()
	rd.JaggerParams = ""

	jsonRequest, _ := json.Marshal(rd)

	// creates pending request to add do the stream pending map
	pr := pending.Request{}
	pr.RequestID = rd.IDRequest
	pr.ResponseChan = make(chan string)
	pr.RequestData = string(jsonRequest)
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
