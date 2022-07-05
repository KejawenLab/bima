package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
)

func Test_Cors(t *testing.T) {
	middleware := Cors{Options: cors.Options{}}

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, 0, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
}
