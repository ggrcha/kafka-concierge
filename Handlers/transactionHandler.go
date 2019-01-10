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

		pr := pending.Request{}
		pr.RequestID = string(letterBytes[rand.Intn(len(letterBytes))])
		pr.ResponseChan = make(chan string)
		defer closeResources(pr)

		go next.ServeHTTP(w, r)

		kafka.NewRequest <- pr

		select {
		case <-time.After(3 * time.Second):
			log.Println(debuggin.Tracer(), "timeout received")
			kafka.ToChan <- pr.RequestID
		case <-pr.ResponseChan:
			log.Println(debuggin.Tracer(), "received response from kafka consumer")
		}
		log.Println(debuggin.Tracer(), "returning ...")

		w.Write([]byte("ok"))
	}
}

func closeResources(pr pending.Request) {
	close(pr.ResponseChan)
}
