package kafka

import (
	debuggin "kernel-concierge/Debuggin"
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

	// starts channel to receive timeouts
	ToChan = make(chan string)
	NewRequest = make(chan pending.Request)

	rand.Seed(time.Now().UnixNano())

	for {

		log.Println(debuggin.Tracer(), "awake")

		select {
		case id := <-ToChan:
			pr := pending.Request{RequestID: id}
			pr.Remove()
		case nr := <-NewRequest:
			log.Println(debuggin.Tracer(), "new request arrived: ", nr)
			nr.Add()
		}

		requestID := string(letterBytes[rand.Intn(len(letterBytes))])
		rp := pending.Request{RequestID: requestID}
		c, exists := rp.GetByID()
		if exists {
			rp.Remove()
			c <- "response received"
		}
	}
}
