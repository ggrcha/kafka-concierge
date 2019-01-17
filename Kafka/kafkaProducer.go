package kafka

import (
	"kernel-concierge/Debuggin"
	"kernel-concierge/Pending"
	"log"
	"os"

	"gopkg.in/Shopify/sarama.v1"
)

var hostname string

// ProduceRequest starts the process of a new transaction
func ProduceRequest(request pending.Request) {

	hostname, _ = os.Hostname()

	producer := getProducer()

	msg := prepareMessage(reqTopic, request.RequestData)

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		log.Println(debuggin.Tracer(), "erro: ", err)
		return
	}
	log.Println(debuggin.Tracer(), "Message persisted.")
	return

}

func prepareMessage(topic string, message interface{}) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message.(string)),
		Key:       sarama.StringEncoder(hostname),
	}

	return msg
}
