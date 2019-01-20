package pending

// StreamPendingRequests keeps all requests pending response
var streamPendingRequests map[string]chan map[string]interface{}

// Request ...
type Request struct {
	RequestID    string
	ResponseChan chan map[string]interface{}
	RequestData  string
}

// Add manages the struct containing pending requests
func (pr Request) Add() {

	// if no map avaiable, creates one
	if streamPendingRequests == nil {
		setStreamPendingRequests(make(map[string]chan map[string]interface{}))
	}

	streamPendingRequests[pr.RequestID] = pr.ResponseChan
}

// Remove ...
func (pr Request) Remove() {
	delete(streamPendingRequests, pr.RequestID)
}

// GetByID ...
func (pr Request) GetByID() (chan map[string]interface{}, bool) {
	r, exists := streamPendingRequests[pr.RequestID]
	return r, exists
}

// SetStreamPendingRequests ...
func setStreamPendingRequests(spr map[string]chan map[string]interface{}) {
	streamPendingRequests = spr
}
