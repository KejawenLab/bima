package brokers

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Kafka struct {
	publisher *kafka.Publisher
	consumer  *kafka.Subscriber
}

func (b *Kafka) Publish(queueName string, payload message.Payload) error {
	return b.publisher.Publish(queueName, message.NewMessage(watermill.NewUUID(), payload))
}

func (b *Kafka) Consume(queueName string) (<-chan *message.Message, error) {
	return b.consumer.Subscribe(context.Background(), queueName)
}
