package routes

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_Hello_Route_Success(t *testing.T) {
	listen, err := net.Listen("tcp", ":111")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	go server.Serve(listen)

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	logger := logrus.New()
	env := configs.Env{
		Debug: true,
	}

	health := Health{
		Logger: &handlers.Logger{
			Env:    &env,
			Logger: logger,
		},
	}

	time.Sleep(100 * time.Millisecond)

	health.SetClient(conn)

	assert.Equal(t, http.MethodGet, health.Method())
	assert.Equal(t, HEALTH_PATH, health.Path())
	assert.Nil(t, health.Middlewares())

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	health.Handle(w, req, map[string]string{})

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	server.Stop()
}

func Test_Hello_Route_Down(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())

	logger := logrus.New()
	env := configs.Env{
		Debug: true,
	}

	route := Health{
		Logger: &handlers.Logger{
			Env:    &env,
			Logger: logger,
		},
	}

	time.Sleep(100 * time.Millisecond)

	route.SetClient(conn)

	assert.Equal(t, http.MethodGet, route.Method())
	assert.Equal(t, HEALTH_PATH, route.Path())
	assert.Nil(t, route.Middlewares())

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	route.Handle(w, req, map[string]string{})

	resp := w.Result()

	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
}
