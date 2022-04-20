package broker

import (
	"app"
	"github.com/Shopify/sarama"
)

func getProducer(config *sarama.Config) sarama.AsyncProducer {

	producer, err := sarama.NewAsyncProducer(app.Yaml.Base.MQ, config)
	if err != nil {
		panic(err)
	}

	return producer
}
