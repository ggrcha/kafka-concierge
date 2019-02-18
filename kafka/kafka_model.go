package kafka

import (
	pending "kernel-concierge/pending"
	"os"
	"sync"

	cluster "github.com/bsm/sarama-cluster"
	"gopkg.in/Shopify/sarama.v1"
)

var (
	kafkaProducer sarama.SyncProducer
	onceP         sync.Once
	onceCg        sync.Once
	err           error
	cg            *cluster.Consumer
	// ToChan timeout channel
	ToChan chan pending.Request
	// NewRequest new request channel
	NewRequest chan pending.Request
	// Cancel ends
	Cancel        chan bool
	zookeeperConn = os.Getenv("ZK_HOST")
	broker        = os.Getenv("KAFKA_HOST")
	port          = os.Getenv("KAFKA_PORT")
	reqTopic      = os.Getenv("RQ_TOPIC")
	rpTopic       = os.Getenv("RP_TOPIC")
)

const (
	cgroup = "consumer-ricko-response"
)
