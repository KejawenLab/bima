package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_RequestIDHeader(t *testing.T) {
	middleware := RequestID{
		Logger: &loggers.Logger{
			Verbose: true,
			Logger:  logrus.New(),
			Data:    logrus.Fields{},
		},
		RequestIDHeader: "X-Request-ID",
	}

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, 259, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
}
