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
	// localRpTopic := rpTopic

	// config := consumergroup.NewConfig()
	config := cluster.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	// config.Consumer.Offsets.ProcessingTimeout = 10 * time.Second
	config.Version = sar.V2_0_0_0
	config.Consumer.Fetch.Max = 5120
	config.Consumer.Fetch.Default = 5120
	config.Consumer.Fetch.Min = 5120

	cg, _ = cluster.NewConsumer([]string{kafkaBroker}, cgroup, []string{localRpTopic}, config)

}

func initConsumer() sarama.PartitionConsumer {

	kafkaBroker := broker + ":" + port
	hostname, _ := os.Hostname()
	localRpTopic := rpTopic + hostname

	consumer, err := sarama.NewConsumer([]string{kafkaBroker}, nil)
	if err != nil {
		panic(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(localRpTopic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	return partitionConsumer
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
