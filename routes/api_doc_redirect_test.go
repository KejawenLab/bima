package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_Api_Doc_Redirect_Route(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	route := ApiDocRedirect{}
	route.SetClient(conn)

	route.Handle(w, req, map[string]string{})

	resp := w.Result()

	assert.Equal(t, http.MethodGet, route.Method())
	assert.Equal(t, ApiDocPath, route.Path())
	assert.Nil(t, route.Middlewares())

	assert.Equal(t, http.StatusPermanentRedirect, resp.StatusCode)
}
