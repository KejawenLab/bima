package brokers

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Amqp struct {
	publisher *amqp.Publisher
	consumer  *amqp.Subscriber
}

func (b *Amqp) Publish(queueName string, payload message.Payload) error {
	return b.publisher.Publish(queueName, message.NewMessage(watermill.NewUUID(), payload))
}

func (b *Amqp) Consume(queueName string) (<-chan *message.Message, error) {
	return b.consumer.Subscribe(context.Background(), queueName)
}
