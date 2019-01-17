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

// RequestData ...
type RequestData struct {
	IDRequest         string      `json:"idRequest"`
	JaggerParams      interface{} `json:"jagerParams"`
	AccountOperations interface{} `json:"accountOperations"`
}

//Ao ...
type Ao struct {
	AccountOperations []struct {
		ID    string `json:"id"`
		Value int    `json:"value"`
	} `json:"accountOperations"`
}

// ResponseData ...
type ResponseData struct {
	ResponseStatus string `json:"status"`
	IDTransaction  string `json:"idTransaction"`
	ResponseTopic  string `json:"responseTopic"`
}

// Manager manages request to kafka
func Manager(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// rest boilerplate
	req, _ := ioutil.ReadAll(r.Body)

	// --------------------------------------

	// PROVISORIO: somente para o Ricko aceitar os requests
	jprms := `{"uber-trace-id":"a817daac97187a30:a817daac97187a30:0:1"}`

	// --------------------------------------

	ao := Ao{}
	_ = json.Unmarshal(req, &ao)
	// creates request that will be send to kafka pipeline
	rd := RequestData{}
	rd.AccountOperations = ao.AccountOperations
	rd.IDRequest = uuid.Must(uuid.NewV4()).String()
	_ = json.Unmarshal([]byte(jprms), &rd.JaggerParams)

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
