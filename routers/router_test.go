package routers

import (
	"context"
	"net/http"
	"testing"

	mocks "github.com/KejawenLab/bima/v4/mocks/routers"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func Test_Router(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router1 := mocks.NewRouter(t)
	router1.On("Handle", ctx, mock.Anything, mock.Anything).Once()
	router1.On("Priority").Return(1).Once()

	router2 := mocks.NewRouter(t)
	router2.On("Handle", ctx, mock.Anything, mock.Anything).Once()
	router2.On("Priority").Return(2).Once()

	factory := Factory{
		Routers: []Router{
			router1,
			router2,
		},
	}

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())

	factory.Sort()
	factory.Handle(ctx, http.NewServeMux(), conn)

	router1.AssertExpectations(t)
	router2.AssertExpectations(t)
}
