package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KejawenLab/bima/v4/loggers"
	"github.com/stretchr/testify/assert"
)

func Test_Hello_Route(t *testing.T) {
	loggers.Default("test")

	time.Sleep(100 * time.Millisecond)

	health := Health{}

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, http.MethodGet, health.Method())
	assert.Equal(t, HelthPath, health.Path())
	assert.Nil(t, health.Middlewares())

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()
	health.Handle(w, req, map[string]string{})

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
