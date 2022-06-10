package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_RequestIDHeader(t *testing.T) {
	middleware := RequestID{
		Logger: &handlers.Logger{
			Verbose: true,
			Logger:  logrus.New(),
		},
		RequestIDHeader: "X-Request-ID",
	}

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, -252, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
}
