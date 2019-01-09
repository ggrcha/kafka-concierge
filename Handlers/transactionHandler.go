package handlers

import (
	pending "kernel-concierge/Pending"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// StreamHandler wraps requests for new transactions
func StreamHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rand.Seed(time.Now().UnixNano())

		log.Printf("Logged connection from %s", r.RemoteAddr)

		go next.ServeHTTP(w, r)

		requestID := letterBytes[rand.Intn(len(letterBytes))]
		responseChan := make(chan bool)

		go pending.ManagePendingRequest(string(requestID), responseChan)

		log.Println(pending.StreamPendingRequests)

		select {
		case <-time.After(1 * time.Second):
			log.Println("timeout received")
		}

		log.Printf("returning ...")

	}
}
