package services

import (
	"encoding/json"
	debuggin "kernel-concierge/Debuggin"
	kafka "kernel-concierge/Kafka"
	pending "kernel-concierge/Pending"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// ASResponseData ...
type ASResponseData struct {
	ResponseStatus string `json:"status"`
	AccountID      string `json:"accountId"`
	ResponseTopic  string `json:"responseTopic"`
}

// RequestAccount ...
type RequestAccount struct {
	RequestID string `json:"requestId"`
}

//AccountService manages new accounts requests
func AccountService(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// creates request that will be send to kafka pipeline
	ra := RequestAccount{}
	ra.RequestID = uuid.Must(uuid.NewV4()).String()
	requestAccount, _ := json.Marshal(ra)

	// creates pending request to add do the stream pending map
	pr := pending.Request{}
	pr.RequestID = ra.RequestID
	pr.ResponseChan = make(chan map[string]interface{})
	pr.RequestData = string(requestAccount)

	// adds newly created request to control map
	kafka.NewRequest <- pr

	// produces message to kafka pipeline
	kafka.ProduceRequest(pr)

	// blocks execution until response or timeout
	select {
	case <-time.After(10 * time.Second):
		log.Println(debuggin.Tracer(), "timeout received")
		// Notifies timeout
		kafka.ToChan <- pr
		rd := ASResponseData{}
		rd.AccountID = ""
		rd.ResponseStatus = "INTERNAL_ERROR"
		rd.ResponseTopic = ""
		resp, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
	case rp := <-pr.ResponseChan:
		// returns response to client
		log.Println(debuggin.Tracer(), "received response from kafka consumer")
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(rp)
		w.Write(resp)
	}
}
