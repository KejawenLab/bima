package routers

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/KejawenLab/bima/v3/configs"
	mocks "github.com/KejawenLab/bima/v3/mocks/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func Test_Gateway_Router(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())

	server := runtime.NewServeMux()

	grpc := mocks.NewServer(t)
	grpc.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	router := GRpcGateway{}
	router.Register([]configs.Server{grpc})

	assert.Equal(t, 255, router.Priority())
	assert.Equal(t, 1, len(router.servers))

	router.Handle(context.TODO(), server, conn)

	req := httptest.NewRequest("GET", "http://bima.framework/handle", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	grpc.AssertExpectations(t)
}
