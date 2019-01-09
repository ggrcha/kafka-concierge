package handlers

import (
	pending "kernel-concierge/Pending"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "abcdefABCDEF"

// StreamHandler wraps requests for new transactions
func StreamHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rand.Seed(time.Now().UnixNano())

		log.Printf("Logged connection from %s", r.RemoteAddr)

		go next.ServeHTTP(w, r)

		requestID := letterBytes[rand.Intn(len(letterBytes))]
		responseChan := make(chan bool)

		go pending.ManagePendingRequest(string(requestID), responseChan)

		select {
		case <-time.After(40 * time.Second):
			log.Println("timeout received")
		case <-responseChan:
			log.Println("received response from kafka consumer")
		}

		log.Printf("returning ...")

	}
}
