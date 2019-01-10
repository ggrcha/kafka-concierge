package kafka

import (
	pending "kernel-concierge/Pending"
	"log"
	"math/rand"
	"time"
)

// ToChan ...
var ToChan chan string

// NewRequest ...
var NewRequest chan pending.Request

const letterBytes = "ab"

// ConsumeKafkaResponses ...
func ConsumeKafkaResponses() {

	// creates the pending requests map and injects it on the pending request routine
	streamPendingRequests := make(map[string]chan string)
	pending.SetStreamPendingRequests(streamPendingRequests)

	// starts channel to receive timeouts
	ToChan = make(chan string)
	NewRequest = make(chan pending.Request)

	rand.Seed(time.Now().UnixNano())

	for {

		log.Println("awake")

		select {
		case id := <-ToChan:
			pr := pending.Request{RequestID: id}
			pr.Remove()
		case nr := <-NewRequest:
			log.Println("new request arrived: ", nr)
			nr.Add()
		}

		requestID := string(letterBytes[rand.Intn(len(letterBytes))])
		rp := pending.Request{RequestID: requestID}
		c, exists := rp.GetByID()
		if exists {
			c <- "response received"
		}
	}
}
