package kafka

import (
	"log"
	"math/rand"
	"time"
)

const letterBytes = "abcdefABCDEF"

// StreamPendingRequests keeps all requests pending response
var StreamPendingRequests map[string]chan bool

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
			log.Println("returning to client")
			rChan <- true
		}
		log.Println(StreamPendingRequests)
	}
}
