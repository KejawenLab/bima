package events

import (
	"errors"
	"testing"

	mocks "github.com/KejawenLab/bima/v4/mocks/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Dispatcher(t *testing.T) {
	listener1 := mocks.NewListener(t)
	listener1.On("Priority").Return(0).Maybe()
	listener1.On("Listen").Return("test_other").Maybe()
	listener1.On("Handle", mock.Anything).Return(mock.Anything).Maybe()

	listener2 := mocks.NewListener(t)
	listener2.On("Priority").Return(0).Maybe()
	listener2.On("Listen").Return(PaginationEvent.String()).Maybe()
	listener2.On("Handle", mock.Anything).Return(mock.Anything).Maybe()

	v := errors.New("")
	dispatcher := Dispatcher{}
	dispatcher.Register([]Listener{listener1, listener2})
	assert.Error(t, dispatcher.Dispatch("unregistered", v))
	assert.Nil(t, dispatcher.Dispatch(PaginationEvent.String(), v))

	listener1.AssertExpectations(t)
	listener2.AssertExpectations(t)
}
