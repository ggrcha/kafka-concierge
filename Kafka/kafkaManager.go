package kafka

import (
	debuggin "kernel-concierge/Debuggin"
	"log"
	"os"
	"time"

	sar "github.com/Shopify/sarama"
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

func getConsumer() {

	kafkaBroker := broker + ":" + port

	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = 10 * time.Second
	config.Version = sar.V2_0_0_0

	conf := cluster.NewConfig()

	hostname, _ := os.Hostname()
	localRpTopic := rpTopic + hostname

	cg, _ = cluster.NewConsumer([]string{kafkaBroker}, cgroup, []string{localRpTopic}, conf)

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
