package messengers

import (
	"errors"
	"testing"

	"github.com/KejawenLab/bima/v4/loggers"
	mocks "github.com/KejawenLab/bima/v4/mocks/messengers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Messenger_No_Broker(t *testing.T) {
	loggers.Default("test")

	topic := "test"
	data := []byte("test")

	messenger := New(true, nil)
	err := messenger.Publish(topic, data)
	assert.Error(t, err)

	_, err = messenger.Consume(topic)
	assert.Error(t, err)
}

func Test_Messenger_Publish_Debug(t *testing.T) {
	loggers.Default("test")

	topic := "test"
	data := []byte("test")

	broker := mocks.NewBroker(t)
	broker.On("Publish", topic, mock.Anything).Return(errors.New("")).Once()

	messenger := New(true, broker)
	assert.Error(t, messenger.Publish(topic, data))

	broker.AssertExpectations(t)

	broker = mocks.NewBroker(t)
	broker.On("Publish", topic, mock.Anything).Return(nil).Once()

	messenger = New(true, broker)
	assert.Nil(t, messenger.Publish(topic, data))

	broker.AssertExpectations(t)
}

func Test_Messenger_Publish_Non_Debug(t *testing.T) {
	loggers.Default("test")

	topic := "test"
	data := []byte("test")

	broker := mocks.NewBroker(t)
	broker.On("Publish", topic, mock.Anything).Return(errors.New("")).Once()

	messenger := New(false, broker)
	assert.Error(t, messenger.Publish(topic, data))

	broker.AssertExpectations(t)

	broker = mocks.NewBroker(t)
	broker.On("Publish", topic, mock.Anything).Return(nil).Once()

	messenger = New(false, broker)
	assert.Nil(t, messenger.Publish(topic, data))

	broker.AssertExpectations(t)
}

func Test_Messenger_Consume_Debug(t *testing.T) {
	loggers.Default("test")

	topic := "test"

	broker := mocks.NewBroker(t)
	broker.On("Consume", topic).Return(nil, errors.New("")).Once()

	messenger := New(true, broker)
	_, err := messenger.Consume(topic)

	assert.Error(t, err)

	broker.AssertExpectations(t)

	broker = mocks.NewBroker(t)
	broker.On("Consume", topic).Return(nil, nil).Once()

	messenger = New(true, broker)
	_, err = messenger.Consume(topic)

	assert.Nil(t, err)

	broker.AssertExpectations(t)
}

func Test_Messenger_Consume_Non_Debug(t *testing.T) {
	loggers.Default("test")

	topic := "test"

	broker := mocks.NewBroker(t)
	broker.On("Consume", topic).Return(nil, errors.New("")).Once()

	messenger := New(false, broker)
	_, err := messenger.Consume(topic)

	assert.Error(t, err)

	broker.AssertExpectations(t)

	broker = mocks.NewBroker(t)
	broker.On("Consume", topic).Return(nil, nil).Once()

	messenger = New(false, broker)
	_, err = messenger.Consume(topic)

	assert.Nil(t, err)

	broker.AssertExpectations(t)
}
