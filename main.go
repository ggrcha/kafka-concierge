package main

import (
	debuggin "kernel-concierge/Debuggin"
	handlers "kernel-concierge/Handlers"
	kafka "kernel-concierge/Kafka"
	pending "kernel-concierge/Pending"
	tran "kernel-concierge/Transaction"
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

	// starts routine that gets kafka's responses
	go kafka.ConsumeKafkaResponses()

	go func() {
		select {
		case <-signals:
			signal.Stop(signals)
			log.Println(debuggin.Tracer(), "exiting gracefully...")
			kafka.Cancel <- true
			closeResources()
		}
	}()

	http.Handle("/log", handlers.StreamHandler(tran.Manager))
	http.ListenAndServe(":8000", nil)
}

func closeResources() {
	close(kafka.ToChan)
	close(kafka.NewRequest)
	close(kafka.Cancel)
}
