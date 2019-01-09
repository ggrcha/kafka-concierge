package main

import (
	handlers "kernel-concierge/Handlers"
	"kernel-concierge/Pending"
	tran "kernel-concierge/Transaction"
	"net/http"
)

func main() {

	StreamPendingRequests := make(map[string]chan bool)
	pending.StreamPendingRequests = StreamPendingRequests

	http.Handle("/log", handlers.StreamHandler(tran.Manager))
	http.ListenAndServe(":8000", nil)
}
