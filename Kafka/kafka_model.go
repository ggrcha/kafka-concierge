package kafka

import (
	pending "kernel-concierge/Pending"
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
	Cancel chan bool
)

const (
	zookeeperConn = "localhost:2181"
	cgroup        = "kernel-concierge"
	reqTopic      = "kernel-concierge-rq"
	rpTopic       = "kernel-concierge-rp"
	broker        = "localhost:9092"
)
