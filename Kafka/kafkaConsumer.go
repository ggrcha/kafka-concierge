package kafka

import (
	"encoding/json"
	debuggin "kernel-concierge/Debuggin"
	pending "kernel-concierge/Pending"
	"log"
	"math/rand"
	"os"
	"time"
)

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
			var rv map[string]interface{}
			if err := json.Unmarshal(msg.Value, &rv); err != nil {
				panic(err)
			}
			requestID := rv["idRequest"].(string)
			rp := pending.Request{RequestID: requestID}
			// gets response channel to send response
			rc, exists := rp.GetByID()
			if exists {
				rp.Remove()
				rc <- string(msg.Value)
				close(rc)
			}
		case pr := <-ToChan:
			pr.Remove()
			close(pr.ResponseChan)
		case nr := <-NewRequest:
			log.Println(debuggin.Tracer(), "new request arrived: ", nr)
			nr.Add()
		case <-Cancel:
			log.Println(debuggin.Tracer(), "kafka exiting")
			closeResources()
			os.Exit(0)
		}
	}
}
