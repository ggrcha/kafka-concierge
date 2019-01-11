package transaction

import (
	"encoding/json"
	"io/ioutil"
	debuggin "kernel-concierge/Debuggin"
	kafka "kernel-concierge/Kafka"
	pending "kernel-concierge/Pending"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// IncomingTransaction ...
type IncomingTransaction struct {
	Message string `json:"message"`
}

// ResponseMessage ...
type ResponseMessage struct {
	Msg string `json:"message"`
}

const letterBytes = "ab"

var rm ResponseMessage

// Manager manages request to kafka
func Manager(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	rm = ResponseMessage{}

	rand.Seed(time.Now().UnixNano())

	// rest boilerplate
	it := IncomingTransaction{}
	req, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(req), &it)

	if err != nil {
		log.Println(debuggin.Tracer(), "invalid transaction")
		rm.Msg = "Invalid transaction format"
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(rm)
		w.Write(resp)
		return
	}

	// creates pending request to add do the stream pending map
	pr := pending.Request{}
	pr.RequestID = string(letterBytes[rand.Intn(len(letterBytes))])
	pr.ResponseChan = make(chan string)
	pr.RequestData = it.Message
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
