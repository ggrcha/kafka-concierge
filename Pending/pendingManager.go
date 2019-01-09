package pending

// StreamPendingRequests keeps all requests pending response
var StreamPendingRequests map[string]chan bool

// ManagePendingRequest manages the struct containing pending requests
func ManagePendingRequest(requestID string, responseChan chan bool) {

	StreamPendingRequests[requestID] = responseChan

}
