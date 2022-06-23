package routers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KejawenLab/bima/v3/middlewares"
	middlewareMocks "github.com/KejawenLab/bima/v3/mocks/middlewares"
	routeMocks "github.com/KejawenLab/bima/v3/mocks/routes"
	"github.com/KejawenLab/bima/v3/routes"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func Test_Mux_Router(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())

	server := runtime.NewServeMux()

	route := routeMocks.NewRoute(t)
	route.On("Path").Return("/without-middleware").Once()
	route.On("Method").Return(http.MethodGet).Once()
	route.On("SetClient", mock.Anything).Once()
	route.On("Middlewares").Return(nil).Once()
	route.On("Handle", mock.Anything, mock.Anything, mock.Anything).Once()

	router := MuxRouter{}
	router.Register([]routes.Route{route})

	assert.Equal(t, -255, router.Priority())
	assert.Equal(t, 1, len(router.routes))

	router.Handle(context.TODO(), server, conn)

	req := httptest.NewRequest("GET", "http://bima.framework/without-middleware", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	route.AssertExpectations(t)

	middleware := middlewareMocks.NewMiddleware(t)
	middleware.On("Attach", mock.Anything, mock.Anything).Return(false).Once()

	route = routeMocks.NewRoute(t)
	route.On("Path").Return("/middleware").Once()
	route.On("Method").Return(http.MethodGet).Once()
	route.On("SetClient", mock.Anything).Once()
	route.On("Middlewares").Return([]middlewares.Middleware{middleware}).Once()
	route.On("Handle", mock.Anything, mock.Anything, mock.Anything).Once()

	router = MuxRouter{}
	router.Register([]routes.Route{route})

	assert.Equal(t, -255, router.Priority())
	assert.Equal(t, 1, len(router.routes))

	router.Handle(context.TODO(), server, conn)

	req = httptest.NewRequest("GET", "http://bima.framework/middleware", nil)
	w = httptest.NewRecorder()

	server.ServeHTTP(w, req)

	route.AssertExpectations(t)

	middleware = middlewareMocks.NewMiddleware(t)
	middleware.On("Attach", mock.Anything, mock.Anything).Return(true).Once()

	route = routeMocks.NewRoute(t)
	route.On("Path").Return("/middleware-stop").Once()
	route.On("Method").Return(http.MethodGet).Once()
	route.On("SetClient", mock.Anything).Once()
	route.On("Middlewares").Return([]middlewares.Middleware{middleware}).Once()

	router = MuxRouter{}
	router.Register([]routes.Route{route})

	assert.Equal(t, -255, router.Priority())
	assert.Equal(t, 1, len(router.routes))

	router.Handle(context.TODO(), server, conn)

	req = httptest.NewRequest("GET", "http://bima.framework/middleware-stop", nil)
	w = httptest.NewRecorder()

	server.ServeHTTP(w, req)

	route.AssertExpectations(t)
}
