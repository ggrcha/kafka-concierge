package main

import (
	handlers "kernel-concierge/Handlers"
	kafka "kernel-concierge/Kafka"
	tran "kernel-concierge/Transaction"
	"net/http"
)

func main() {

	go kafka.ConsumeKafkaResponses()

	http.Handle("/log", handlers.StreamHandler(tran.Manager))
	http.ListenAndServe(":8000", nil)
}
