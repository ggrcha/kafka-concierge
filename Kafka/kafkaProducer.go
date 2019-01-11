package kafka

import (
	"kernel-concierge/Debuggin"
	"kernel-concierge/Pending"
	"log"

	"gopkg.in/Shopify/sarama.v1"
)

var brokers = []string{"127.0.0.1:9092"}
var topic = "kernel-concierge"

// ProduceRequest starts the process of a new transaction
func ProduceRequest(request pending.Request) {

	producer := getProducer()

	msg := prepareMessage(topic, request.RequestData)
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		log.Println(debuggin.Tracer(), "erro: ", err)
		return
	}
	log.Println(debuggin.Tracer(), "Message persisted.")
	return

}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}
