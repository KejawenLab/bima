package messengers

import (
	"context"
	"errors"
	"strings"

	"github.com/KejawenLab/bima/v4/loggers"
	"github.com/ThreeDotsLabs/watermill/message"
)

type (
	Broker interface {
		Publish(queueName string, payload message.Payload) error
		Consume(queueName string) (<-chan *message.Message, error)
	}

	Messenger struct {
		debug  bool
		broker Broker
	}
)

func New(debug bool, broker Broker) *Messenger {
	return &Messenger{
		debug:  debug,
		broker: broker,
	}
}

func (m *Messenger) Publish(queueName string, data []byte) error {
	ctx := context.WithValue(context.Background(), "scope", "messenger")
	if m.broker == nil {
		loggers.Logger.Error(ctx, "broker not configured properly")

		return errors.New("broker not configured properly")
	}

	if m.debug {
		var log strings.Builder
		log.WriteString("publishing message to: ")
		log.WriteString(queueName)

		loggers.Logger.Debug(ctx, log.String())
	}

	if err := m.broker.Publish(queueName, data); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return err
	}

	return nil
}

func (m *Messenger) Consume(queueName string) (<-chan *message.Message, error) {
	ctx := context.WithValue(context.Background(), "scope", "messenger")
	if m.broker == nil {
		loggers.Logger.Error(ctx, "broker not configured properly")

		return nil, errors.New("broker not configured properly")
	}

	if m.debug {
		var log strings.Builder
		log.WriteString("consuming: ")
		log.WriteString(queueName)

		loggers.Logger.Debug(ctx, log.String())
	}

	messages, err := m.broker.Consume(queueName)
	if err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, err
	}

	return messages, nil
}
