package handlers

import (
	"context"
	"fmt"

	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Messenger struct {
	Publisher *amqp.Publisher
	Consumer  *amqp.Subscriber
	Logger    *loggers.Logger
}

func (m *Messenger) Publish(queueName string, data []byte) error {
	ctx := context.WithValue(context.Background(), "scope", "messenger")
	m.Logger.Debug(ctx, fmt.Sprintf("Publishing message to: %s", queueName))

	msg := message.NewMessage(watermill.NewUUID(), data)
	if err := m.Publisher.Publish(queueName, msg); err != nil {
		m.Logger.Error(ctx, err.Error())

		return err
	}

	return nil
}

func (m *Messenger) Consume(queueName string) (<-chan *message.Message, error) {
	ctx := context.WithValue(context.Background(), "scope", "messenger")
	m.Logger.Debug(ctx, fmt.Sprintf("Consuming: %s", queueName))

	messages, err := m.Consumer.Subscribe(context.Background(), queueName)
	if err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, err
	}

	return messages, nil
}
