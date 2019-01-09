package kafka

import (
	"log"
	"math/rand"
	"time"
)

// StreamPendingRequests keeps all requests pending response
var StreamPendingRequests map[string]chan bool

const letterBytes = "abcdefABCDEF"

// ConsumeKafkaResponses ...
func ConsumeKafkaResponses() {

	rand.Seed(time.Now().UnixNano())
	var requestID string

	for {
		time.Sleep(1 * time.Second)

		requestID = string(letterBytes[rand.Intn(len(letterBytes))])
		rChan, exists := StreamPendingRequests[requestID]
		log.Println("exists? ", exists)
		if exists {
			log.Println("returning to send response to client")
			rChan <- true
			delete(StreamPendingRequests, requestID)
		}
		log.Println(StreamPendingRequests)
	}
}
