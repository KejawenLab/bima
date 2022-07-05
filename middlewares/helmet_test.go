package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Helmet(t *testing.T) {
	middleware := Helmet{}

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, -255, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
}
