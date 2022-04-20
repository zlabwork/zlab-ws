package business

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"strings"
)

func consumer(config *sarama.Config) {

	address := strings.Split(os.Getenv("KAFKA_ADDRS"), ",")
	consumer, err := sarama.NewConsumer(address, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("Consumer started")

	partitionConsumer, err := consumer.ConsumePartition("message", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	consumed := 0
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			log.Println(msg.Value)
			consumed++

		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
