package business

import (
	"app"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

func consumer(ch chan *[]byte, config *sarama.Config) {

	consumer, err := sarama.NewConsumer(app.Yaml.MQ, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

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
			// log.Printf("Consumed message offset %d\n", msg.Offset)
			ch <- &msg.Value
			consumed++

		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
