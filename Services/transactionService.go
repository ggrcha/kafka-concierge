package services

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
	JaegerParams      interface{} `json:"jaegerParams"`
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

// JP ...
type JP struct {
	UberTraceID interface{} `json:"uber-trace-id"`
}

// TransactionService manages transaction requests to kafka
func TransactionService(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// rest boilerplate
	req, _ := ioutil.ReadAll(r.Body)

	// --------------------------------------

	// PROVISORIO: somente para o Ricko aceitar os requests
	// tracer, closer := jaeger.InitJaeger("hello-world")
	// defer closer.Close()

	// span := tracer.StartSpan("say-hello")
	// log.Println("ctx: ", span.Context())

	jp := JP{}
	jp.UberTraceID = "2648669eebc21e7a:2648669eebc21e7a:0:1"
	// jp.UberTraceID = fmt.Sprint(span)

	// defer span.Finish()

	// --------------------------------------

	ao := Ao{}
	_ = json.Unmarshal(req, &ao)
	// creates request that will be send to kafka pipeline
	rd := RequestData{}
	rd.AccountOperations = ao.AccountOperations
	rd.IDRequest = uuid.Must(uuid.NewV4()).String()
	jParms, _ := json.Marshal(jp)
	_ = json.Unmarshal(jParms, &rd.JaegerParams)

	// converts to json
	jsonRequest, _ := json.Marshal(rd)

	// creates pending request to add do the stream pending map
	pr := pending.Request{}
	pr.RequestID = rd.IDRequest
	pr.ResponseChan = make(chan map[string]interface{})
	pr.RequestData = string(jsonRequest)

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
