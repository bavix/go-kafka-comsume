package main

import (
	"context"
	"errors"
	"github.com/bavix/go-kafka-consume/avro"
	"github.com/bavix/go-kafka-consume/configure"
	"github.com/bavix/go-kafka-consume/workers/retailrocket"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/riferrei/srclient"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(errors.New("no .env file found"))
	}
}

func main() {
	config := configure.New()
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(config.SchemaRegistry.URL)
	marshal := avro.AvroMarshaler{
		SchemaRegistryClient: schemaRegistryClient,
	}

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               config.BrokerList,
			Unmarshaler:           marshal,
			OverwriteSaramaConfig: config.Sarama,
			ConsumerGroup:         config.Consumer.GroupID,
		},
		watermill.NewStdLogger(true, true),
	)

	if err != nil {
		panic(err)
	}

	messages, err := subscriber.Subscribe(context.Background(), config.RetailRocket.Topic)
	if err != nil {
		panic(err)
	}

	retailRocket := retailrocket.RetailRocket{
		HttpClient: retailrocket.NewClient(config.RetailRocket),
	}

	retailRocket.ReadMessages(messages)
}
