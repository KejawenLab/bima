package brokers

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
)

type GoogleCloud struct {
	publisher *googlecloud.Publisher
	consumer  *googlecloud.Subscriber
}

func (b *GoogleCloud) Publish(queueName string, payload message.Payload) error {
	return b.publisher.Publish(queueName, message.NewMessage(watermill.NewUUID(), payload))
}

func (b *GoogleCloud) Consume(queueName string) (<-chan *message.Message, error) {
	return b.consumer.Subscribe(context.Background(), queueName)
}
