package main

import (
	handlers "kernel-concierge/Handlers"
	kafka "kernel-concierge/Kafka"
	"kernel-concierge/Pending"
	tran "kernel-concierge/Transaction"
	"net/http"
)

func main() {

	// Creates the pending requests map and injects it on the pending request routine
	StreamPendingRequests := make(map[string]chan bool)
	pending.StreamPendingRequests = StreamPendingRequests
	kafka.StreamPendingRequests = StreamPendingRequests

	go kafka.ConsumeKafkaResponses()

	http.Handle("/log", handlers.StreamHandler(tran.Manager))
	http.ListenAndServe(":8000", nil)
}
