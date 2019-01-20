package kafka

import (
	"encoding/json"
	debuggin "kernel-concierge/Debuggin"
	pending "kernel-concierge/Pending"
	"log"
	"os"
)

// ConsumeKafkaResponses ...
func ConsumeKafkaResponses() {

	cg := getConsumer()

	for {

		log.Println(debuggin.Tracer(), "awake")

		// blocks execution until some timeout, arrival of new request or new message
		select {
		case msg := <-cg.Messages():
			log.Println(debuggin.Tracer(), "received message: ", string(msg.Key), string(msg.Value))
			// retrieves idRequest from kafka response
			var rv map[string]interface{}
			if err := json.Unmarshal(msg.Value, &rv); err != nil {
				panic(err)
			}
			// requestID := rv["idRequest"].(string)
			requestID := string(msg.Key)
			rp := pending.Request{RequestID: requestID}
			// gets response channel to send response
			rc, exists := rp.GetByID()
			if exists {
				rp.Remove()
				rc <- rv
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
