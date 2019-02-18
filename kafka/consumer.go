package kafka

import (
	"encoding/json"
	debuggin "kernel-concierge/debuggin"
	pending "kernel-concierge/pending"
	"log"
	"os"
)

// ConsumeKafkaResponses ...
func ConsumeKafkaResponses() {

	initConsumerGroup()

	// consume errors
	go func() {
		for err := range cg.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range cg.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	for {

		log.Println(debuggin.Tracer(), "awake")

		// blocks execution until some timeout, arrival of new request or new message
		select {
		case msg := <-cg.Messages():
			log.Println(debuggin.Tracer(), "received message: ", string(msg.Key), string(msg.Value))
			var rv map[string]interface{}
			if err := json.Unmarshal(msg.Value, &rv); err != nil {
				panic(err)
			}
			requestID := string(msg.Key)
			rp := pending.Request{RequestID: requestID}
			// gets response channel to send response
			rc, exists := rp.GetByID()
			if exists {
				rp.Remove()
				rc <- rv
				close(rc)
			}
			cg.MarkOffset(msg, "")
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
