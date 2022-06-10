package middlewares

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Bima_Header(t *testing.T) {
	middleware := Header{}

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, -256, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
}
