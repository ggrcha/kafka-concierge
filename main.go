package main

import (
	_ "encoding/json"
	"fmt"
	_ "io"
	_ "log"
	_ "net/http"
	_ "os"
	_ "sync"
	_ "time"

	_ "github.com/Shopify/sarama"
	_ "github.com/bsm/sarama-cluster"
	_ "github.com/satori/go.uuid"
	_ "github.com/uber/jaeger-client-go"
	_ "github.com/uber/jaeger-client-go/config"
	_ "gopkg.in/Shopify/sarama.v1"
)

func main() {
	fmt.Println("This is used for build caching purposes. Should be replaced.")
}
