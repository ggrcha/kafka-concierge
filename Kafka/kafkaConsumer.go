package kafka

import (
	debuggin "kernel-concierge/Debuggin"
	pending "kernel-concierge/Pending"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var signals chan os.Signal

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
	// Trap SIGINT to trigger a shutdown.
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	cg := getConsumer()

	rand.Seed(time.Now().UnixNano())

	for {

		log.Println(debuggin.Tracer(), "awake")

		// blocks execution until some timeout, arrival of new request or new message
		select {
		case msg := <-cg.Messages():
			log.Println(debuggin.Tracer(), "received message: ", string(msg.Value))
		case pr := <-ToChan:
			pr.Remove()
		case nr := <-NewRequest:
			log.Println(debuggin.Tracer(), "new request arrived: ", nr)
			nr.Add()
		case <-signals:
			log.Println(debuggin.Tracer(), "exiting gracefully...")
			closeResources()
			os.Exit(0)
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
