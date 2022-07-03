package messengers

import (
	"context"
	"errors"
	"strings"

	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Messenger struct {
	debug     bool
	publisher *amqp.Publisher
	consumer  *amqp.Subscriber
}

func New(debug bool, publisher *amqp.Publisher, consumer *amqp.Subscriber) *Messenger {
	return &Messenger{
		debug:     debug,
		publisher: publisher,
		consumer:  consumer,
	}
}

func (m *Messenger) Publish(queueName string, data []byte) error {
	ctx := context.WithValue(context.Background(), "scope", "messenger")
	if m.publisher == nil {
		loggers.Logger.Fatal(ctx, "publisher not configured properly")

		return errors.New("publisher not configured properly")
	}

	if m.debug {
		var log strings.Builder
		log.WriteString("publishing message to: ")
		log.WriteString(queueName)

		loggers.Logger.Debug(ctx, log.String())
	}

	msg := message.NewMessage(watermill.NewUUID(), data)
	if err := m.publisher.Publish(queueName, msg); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return err
	}

	return nil
}

func (m *Messenger) Consume(queueName string) (<-chan *message.Message, error) {
	ctx := context.WithValue(context.Background(), "scope", "messenger")
	if m.consumer == nil {
		loggers.Logger.Fatal(ctx, "consumer not configured properly")

		return nil, errors.New("consumer not configured properly")
	}

	if m.debug {
		var log strings.Builder
		log.WriteString("consuming: ")
		log.WriteString(queueName)

		loggers.Logger.Debug(ctx, log.String())
	}

	messages, err := m.consumer.Subscribe(context.Background(), queueName)
	if err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, err
	}

	return messages, nil
}
