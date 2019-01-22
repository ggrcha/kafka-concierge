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

// AccountRequest ...
type AccountRequest struct {
	AccountID string `json:"accountId"`
}

//AccountService manages new accounts requests
func AccountService(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// creates request that will be send to kafka pipeline
	ar := AccountRequest{}
	ar.AccountID = uuid.Must(uuid.NewV4()).String()
	accountRequest, _ := json.Marshal(ar)

	log.Println(debuggin.Tracer(), "accountRequest: ", string(accountRequest))
	// creates pending request to add do the stream pending map
	pr := pending.Request{}
	pr.RequestID = ar.AccountID
	pr.ResponseChan = make(chan map[string]interface{})
	pr.RequestData = string(accountRequest)

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
		rd := ResponseData{}
		rd.IDTransaction = pr.RequestID
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
