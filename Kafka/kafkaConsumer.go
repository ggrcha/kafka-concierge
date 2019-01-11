package kafka

import (
	debuggin "kernel-concierge/Debuggin"
	pending "kernel-concierge/Pending"
	"log"
	"math/rand"
	"time"
)

// ToChan timeout channel
var ToChan chan pending.Request

// NewRequest new request channel
var NewRequest chan pending.Request

const letterBytes = "ab"

// ConsumeKafkaResponses ...
func ConsumeKafkaResponses() {

	// starts channel to receive timeouts
	ToChan = make(chan pending.Request)
	NewRequest = make(chan pending.Request)

	rand.Seed(time.Now().UnixNano())

	for {

		log.Println(debuggin.Tracer(), "awake")

		// blocks execution until some timeout or arrival of new request
		select {
		case pr := <-ToChan:
			pr.Remove()
		case nr := <-NewRequest:
			log.Println(debuggin.Tracer(), "new request arrived: ", nr)
			nr.Add()
		}

		// receives new response from kafka
		requestID := string(letterBytes[rand.Intn(len(letterBytes))])
		rp := pending.Request{RequestID: requestID}
		// gets response channel to send response
		c, exists := rp.GetByID()
		if exists {
			rp.Remove()
			c <- "response received"
		}
	}
}
