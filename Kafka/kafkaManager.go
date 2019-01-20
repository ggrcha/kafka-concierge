package kafka

import (
	debuggin "kernel-concierge/Debuggin"
	"log"
	"os"
	"time"

	"github.com/bsm/sarama-cluster"
	"github.com/wvanbergen/kafka/consumergroup"
	"gopkg.in/Shopify/sarama.v1"
)

func getProducer() sarama.SyncProducer {
	onceP.Do(func() {
		kafkaProducer, err = newProducer()
		if err != nil {
			log.Println(debuggin.Tracer(), "Could not create producer: ", err)
			panic(err)
		}
	})
	return kafkaProducer
}

// func getConsumer() *consumergroup.ConsumerGroup {
func getConsumer() *cluster.Consumer {
	// onceCg.Do(func() {
	// consumer config

	kafkaBroker := broker + ":" + port

	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetOldest
	config.Offsets.ProcessingTimeout = 10 * time.Second

	conf := cluster.NewConfig()

	hostname, _ := os.Hostname()
	localRpTopic := rpTopic + hostname

	consumer, _ := cluster.NewConsumer([]string{kafkaBroker}, cgroup, []string{localRpTopic}, conf)

	return consumer
}

func newProducer() (sarama.SyncProducer, error) {

	kafkaBroker := broker + ":" + port

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Net.SASL.Enable = false
	config.Net.SASL.Handshake = false
	config.Net.TLS.Enable = false
	config.Version = sarama.V2_0_0_0
	producer, err := sarama.NewSyncProducer([]string{kafkaBroker}, config)

	return producer, err
}

// CloseResources closes all resources
func closeResources() {
	kafkaProducer.Close()
	// cg.Close()
}
