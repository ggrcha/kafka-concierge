package kafka

import (
	debuggin "kernel-concierge/Debuggin"
	"log"
	"os"

	sar "github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
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

func initConsumerGroup() {

	kafkaBroker := broker + ":" + port
	hostname, _ := os.Hostname()
	localRpTopic := rpTopic + hostname

	config := cluster.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sar.V2_0_0_0
	config.Consumer.Fetch.Max = 5 * 1024
	config.Consumer.Fetch.Default = 5 * 1024
	config.Consumer.Fetch.Min = 1

	cg, err = cluster.NewConsumer([]string{kafkaBroker}, cgroup, []string{localRpTopic}, config)

	if err != nil {
		log.Println(debuggin.Tracer(), err)
		panic(err)
	}

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
	cg.Close()
}
