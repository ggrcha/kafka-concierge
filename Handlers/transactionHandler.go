package handlers

import (
	debuggin "kernel-concierge/Debuggin"
	kafka "kernel-concierge/Kafka"
	pending "kernel-concierge/Pending"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "ab"

// StreamHandler wraps requests for new transactions
func StreamHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rand.Seed(time.Now().UnixNano())

		log.Println(debuggin.Tracer(), "Logged connection from: ", r.RemoteAddr)

		// creates pending request to add do the stream pending map
		pr := pending.Request{}
		pr.RequestID = string(letterBytes[rand.Intn(len(letterBytes))])
		pr.ResponseChan = make(chan string)
		// defers closing the channel and any other resource opened
		defer closeResources(pr)

		// adds newly created request to control map
		kafka.NewRequest <- pr

		// publishs reques to kafka
		go next.ServeHTTP(w, r)

		// blocks execution until response or timeout
		select {
		case <-time.After(3 * time.Second):
			log.Println(debuggin.Tracer(), "timeout received")
			kafka.ToChan <- pr.RequestID
		case <-pr.ResponseChan:
			log.Println(debuggin.Tracer(), "received response from kafka consumer")
		}
		w.Write([]byte("ok"))
	}
}

func closeResources(pr pending.Request) {
	close(pr.ResponseChan)
}
