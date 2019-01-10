package pending

import (
	debuggin "kernel-concierge/Debuggin"
	"log"
)

// StreamPendingRequests keeps all requests pending response
var streamPendingRequests map[string]chan string

// Request ...
type Request struct {
	RequestID    string
	ResponseChan chan string
}

// Add manages the struct containing pending requests
func (pr Request) Add() {
	streamPendingRequests[pr.RequestID] = pr.ResponseChan
	log.Println(debuggin.Tracer(), "NewRequest: ", streamPendingRequests)
}

// Remove ...
func (pr Request) Remove() {
	delete(streamPendingRequests, pr.RequestID)
	log.Println(debuggin.Tracer(), "ToChan: ", streamPendingRequests)
}

// GetByID ...
func (pr Request) GetByID() (chan string, bool) {
	r, exists := streamPendingRequests[pr.RequestID]
	return r, exists
}

// SetStreamPendingRequests ...
func SetStreamPendingRequests(spr map[string]chan string) {
	streamPendingRequests = spr
}
