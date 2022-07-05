package routes

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KejawenLab/bima/v4/loggers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_Hello_Route_Success(t *testing.T) {
	loggers.Default("test")

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

	health := Health{}

	time.Sleep(100 * time.Millisecond)

	health.SetClient(conn)

	assert.Equal(t, http.MethodGet, health.Method())
	assert.Equal(t, HelthPath, health.Path())
	assert.Nil(t, health.Middlewares())

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
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

	health := Health{}

	time.Sleep(100 * time.Millisecond)

	health.SetClient(conn)

	assert.Equal(t, http.MethodGet, health.Method())
	assert.Equal(t, HelthPath, health.Path())
	assert.Nil(t, health.Middlewares())

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()
	health.Handle(w, req, map[string]string{})

	resp := w.Result()

	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
}
