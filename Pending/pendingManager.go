package pending

// StreamPendingRequests keeps all requests pending response
var StreamPendingRequests map[string]chan bool

// AddPendingRequest manages the struct containing pending requests
func AddPendingRequest(requestID string, responseChan chan bool) {

	StreamPendingRequests[requestID] = responseChan

}

// RemovePendingRequest ...
func RemovePendingRequest(requestID string) {
	delete(StreamPendingRequests, requestID)
}
