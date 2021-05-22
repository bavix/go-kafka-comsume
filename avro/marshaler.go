package avro

import (
	"encoding/binary"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/riferrei/srclient"
)

type AvroMarshaler struct {
	SchemaRegistryClient *srclient.SchemaRegistryClient
}

func (AvroMarshaler) Marshal(topic string, msg *message.Message) (*sarama.ProducerMessage, error) {
	return nil, nil
}

func (m AvroMarshaler) Unmarshal(kafkaMsg *sarama.ConsumerMessage) (*message.Message, error) {
	marshal := kafka.DefaultMarshaler{}
	msg, err := marshal.Unmarshal(kafkaMsg)
	if err != nil {
		return nil, err
	}

	schemaID := binary.BigEndian.Uint32(msg.Payload[1:5])
	schema, err := m.SchemaRegistryClient.GetSchema(int(schemaID))
	if err != nil {
		return nil, err
	}

	codec := schema.Codec()
	native, _, err := codec.NativeFromBinary(msg.Payload[5:])
	if err != nil {
		return nil, err
	}

	msg.Payload, err = codec.TextualFromNative(nil, native)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
