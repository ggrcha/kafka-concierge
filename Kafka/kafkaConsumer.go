package kafka

import (
	debuggin "kernel-concierge/Debuggin"
	pending "kernel-concierge/Pending"
	"log"
	"math/rand"
	"os"
	"time"
)

const letterBytes = "ab"

// ConsumeKafkaResponses ...
func ConsumeKafkaResponses() {

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
		case <-Cancel:
			log.Println(debuggin.Tracer(), "kafka exiting")
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
