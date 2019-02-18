package main

import (
	debuggin "kernel-concierge/debuggin"
	handlers "kernel-concierge/handlers"
	kafka "kernel-concierge/kafka"
	pending "kernel-concierge/pending"
	services "kernel-concierge/services"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	// starts channel to receive timeouts and injects them into kafka
	kafka.ToChan = make(chan pending.Request)
	kafka.NewRequest = make(chan pending.Request)
	kafka.Cancel = make(chan bool)

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		select {
		case <-signals:
			signal.Stop(signals)
			log.Println(debuggin.Tracer(), "exiting gracefully...")
			kafka.Cancel <- true
			closeResources()
		}
	}()

	// starts routine that gets kafka's responses
	go kafka.ConsumeKafkaResponses()

	http.Handle("/v0/transaction", handlers.StreamHandler(services.Transaction))
	http.Handle("/v0/accounts", handlers.StreamHandler(services.Account))
	http.ListenAndServe(":8000", nil)
}

func closeResources() {
	close(kafka.ToChan)
	close(kafka.NewRequest)
	close(kafka.Cancel)
}
