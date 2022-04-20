package broker

import (
	"github.com/Shopify/sarama"
	"os"
	"strings"
)

func getProducer(config *sarama.Config) sarama.AsyncProducer {

	address := strings.Split(os.Getenv("KAFKA_ADDRS"), ",")
	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		panic(err)
	}

	return producer
}
