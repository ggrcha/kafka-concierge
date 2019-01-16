package kafka

import (
	pending "kernel-concierge/Pending"
	"os"
	"sync"

	"github.com/wvanbergen/kafka/consumergroup"
	"gopkg.in/Shopify/sarama.v1"
)

var (
	kafkaProducer sarama.SyncProducer
	onceP         sync.Once
	onceCg        sync.Once
	err           error
	cg            *consumergroup.ConsumerGroup
	// ToChan timeout channel
	ToChan chan pending.Request
	// NewRequest new request channel
	NewRequest chan pending.Request
	// Cancel ends
	Cancel        chan bool
	zookeeperConn = os.Getenv("ZK_HOST")
	broker        = os.Getenv("KAFKA_HOST")
	reqTopic      = os.Getenv("RQ_TOPIC")
	rpTopic       = os.Getenv("RP_TOPIC")
)

const (
	cgroup = "kernel-concierge"
)
