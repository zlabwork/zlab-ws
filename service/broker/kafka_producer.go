package broker

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

func Producer(topic string, config *sarama.Config) {

	address := strings.Split(os.Getenv("KAFKA_ADDRS"), ",")
	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("Producer started")

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// TODO:: bytes data
	bs := []byte("testing 123 " + time.Now().Format(time.RFC3339))
	var enqueued, producerErrors int
ProducerLoop:
	for {
		select {
		case producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(bs)}:
			enqueued++
		case err := <-producer.Errors():
			log.Println("Failed to produce message", err)
			producerErrors++
		case <-signals:
			break ProducerLoop
		}
	}

	log.Printf("Enqueued: %d; errors: %d\n", enqueued, producerErrors)

}
