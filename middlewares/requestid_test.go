package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RequestIDHeader(t *testing.T) {
	middleware := RequestID{
		RequestIDHeader: "X-Request-ID",
	}

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, 259, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
}
