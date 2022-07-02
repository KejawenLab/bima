package routes

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/utils"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_Jwt_Refresh_Invalid_Payload(t *testing.T) {
	loggers.Default("test")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())

	route := JwtRefresh{
		PathUrl:       "/refresh",
		Secret:        "secret",
		SigningMethod: jwt.SigningMethodHS512.Name,
		Expire:        730,
	}
	route.SetClient(conn)

	req := httptest.NewRequest("POST", "http://bima.framework/refresh", nil)
	w := httptest.NewRecorder()
	route.Handle(w, req, map[string]string{})

	resp := w.Result()

	assert.Equal(t, http.MethodPost, route.Method())
	assert.Equal(t, "/refresh", route.Path())
	assert.Nil(t, route.Middlewares())
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload := map[string]string{
		"invalid": "invalid",
	}
	body, _ := json.Marshal(payload)

	req = httptest.NewRequest("POST", "http://bima.framework/refresh", bytes.NewReader(body))
	w = httptest.NewRecorder()
	route.Handle(w, req, map[string]string{})

	resp = w.Result()

	assert.Equal(t, http.MethodPost, route.Method())
	assert.Equal(t, "/refresh", route.Path())
	assert.Nil(t, route.Middlewares())
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = map[string]string{
		"token": "invalid",
	}
	body, _ = json.Marshal(payload)

	req = httptest.NewRequest("POST", "http://bima.framework/refresh", bytes.NewReader(body))
	w = httptest.NewRecorder()
	route.Handle(w, req, map[string]string{})

	resp = w.Result()

	assert.Equal(t, http.MethodPost, route.Method())
	assert.Equal(t, "/refresh", route.Path())
	assert.Nil(t, route.Middlewares())
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_Jwt_Refresh_Valid_Payload(t *testing.T) {
	loggers.Default("test")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())

	route := JwtRefresh{
		PathUrl:       "/refresh",
		Secret:        "secret",
		SigningMethod: jwt.SigningMethodHS512.Name,
		Expire:        730,
	}
	route.SetClient(conn)

	claims := jwt.MapClaims{
		"user": "test",
	}
	token, err := utils.CreateToken("secret", jwt.SigningMethodHS512.Name, claims, 2)
	assert.Nil(t, err)

	refreshToken, err := utils.CreateRefreshToken("secret", jwt.SigningMethodHS512.Name, token)
	assert.Nil(t, err)

	payload := map[string]string{
		"token": refreshToken,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "http://bima.framework/refresh", bytes.NewReader(body))
	w := httptest.NewRecorder()
	route.Handle(w, req, map[string]string{})

	resp := w.Result()

	assert.Equal(t, http.MethodPost, route.Method())
	assert.Equal(t, "/refresh", route.Path())
	assert.Nil(t, route.Middlewares())
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
