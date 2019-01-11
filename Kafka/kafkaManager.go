package kafka

import (
	debuggin "kernel-concierge/Debuggin"
	"log"
	"sync"

	"gopkg.in/Shopify/sarama.v1"
)

var kafkaProducer sarama.SyncProducer
var once sync.Once
var err error

func getProducer() sarama.SyncProducer {
	once.Do(func() {
		kafkaProducer, err = newProducer()
		if err != nil {
			log.Println(debuggin.Tracer(), "Could not create producer: ", err)
			return
		}
	})
	return kafkaProducer
}

func getConsumer() {

}

func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Net.SASL.Enable = false
	config.Net.SASL.Handshake = false
	config.Net.TLS.Enable = false
	producer, err := sarama.NewSyncProducer(brokers, config)

	return producer, err
}
