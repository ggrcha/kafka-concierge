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

// ResponseData ...
type ResponseData struct {
	ResponseStatus string `json:"status"`
	IDTransaction  string `json:"idTransaction"`
	ResponseTopic  string `json:"responseTopic"`
}

var rm ResponseMessage

// Manager manages request to kafka
func Manager(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	rm = ResponseMessage{}

	// rest boilerplate
	req, _ := ioutil.ReadAll(r.Body)

	// creates request that will be send to kafka pipeline
	rd := RequestData{}
	_ = json.Unmarshal(req, &rd.Accounts)
	rd.IDRequest = uuid.Must(uuid.NewV4()).String()
	rd.JaggerParams = ""

	// converts to json
	jsonRequest, _ := json.Marshal(rd)

	// creates pending request to add do the stream pending map
	pr := pending.Request{}
	pr.RequestID = rd.IDRequest
	pr.ResponseChan = make(chan string)
	pr.RequestData = string(jsonRequest)

	// adds newly created request to control map
	kafka.NewRequest <- pr

	// produces message to kafka pipeline
	kafka.ProduceRequest(pr)

	// blocks execution until response or timeout
	select {
	case <-time.After(3 * time.Second):
		log.Println(debuggin.Tracer(), "timeout received")
		// Notifies timeout
		kafka.ToChan <- pr
	case rp := <-pr.ResponseChan:
		// returns response to client
		log.Println(debuggin.Tracer(), "received response from kafka consumer")
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(rp)
		w.Write(resp)
	}
}
