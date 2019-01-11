package kafka

import (
	debuggin "kernel-concierge/Debuggin"
	"log"
	"sync"
	"time"

	"github.com/wvanbergen/kafka/consumergroup"
	"gopkg.in/Shopify/sarama.v1"
)

var (
	kafkaProducer sarama.SyncProducer
	onceP         sync.Once
	onceCg        sync.Once
	err           error
	cg            *consumergroup.ConsumerGroup
)

const (
	zookeeperConn = "127.0.0.1:2181"
	cgroup        = "kernel-concierge"
	reqTopic      = "kernel-concierge-rq"
	rpTopic       = "kernel-concierge-rp"
	broker        = "127.0.0.1:9092"
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

func getConsumer() *consumergroup.ConsumerGroup {

	onceCg.Do(func() {
		// consumer config
		config := consumergroup.NewConfig()
		config.Offsets.Initial = sarama.OffsetOldest
		config.Offsets.ProcessingTimeout = 10 * time.Second

		// join to consumer group
		cg, _ = consumergroup.JoinConsumerGroup(cgroup, []string{rpTopic}, []string{zookeeperConn}, config)
		if err != nil {
			panic(err)
		}
	})

	return cg
}

func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Net.SASL.Enable = false
	config.Net.SASL.Handshake = false
	config.Net.TLS.Enable = false
	producer, err := sarama.NewSyncProducer([]string{broker}, config)

	return producer, err
}

// CloseResources closes all resources
func closeResources() {
	kafkaProducer.Close()
	cg.Close()
}
