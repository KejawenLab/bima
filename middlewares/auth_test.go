package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/stretchr/testify/assert"
)

func Test_Auth(t *testing.T) {
	middleware := Auth{
		Env: &configs.Env{
			AuthHeader: configs.AuthHeader{
				Id:    "",
				Email: "",
				Role:  "",
			},
		},
	}

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	middleware = Auth{
		Env: &configs.Env{
			AuthHeader: configs.AuthHeader{
				Id:        "X-User-Id",
				Email:     "X-User-Email",
				Role:      "X-User-Role",
				Whitelist: "/foo",
			},
		},
	}

	req = httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w = httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	middleware = Auth{
		Env: &configs.Env{
			AuthHeader: configs.AuthHeader{
				Id:        "X-User-Id",
				Email:     "X-User-Email",
				Role:      "X-User-Role",
				Whitelist: "/not-secure",
				MinRole:   2,
			},
		},
	}

	req = httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w = httptest.NewRecorder()

	req.Header.Add("X-User-Id", "1")
	req.Header.Add("X-User-Email", "surya@bima.com")
	req.Header.Add("X-User-Role", "1")

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, true, middleware.Attach(req, w))
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

	middleware = Auth{
		Env: &configs.Env{
			AuthHeader: configs.AuthHeader{
				Id:        "X-User-Id",
				Email:     "X-User-Email",
				Role:      "X-User-Role",
				Whitelist: "/not-secure",
				MinRole:   2,
			},
		},
	}

	req = httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w = httptest.NewRecorder()

	req.Header.Add("X-User-Id", "1")
	req.Header.Add("X-User-Email", "surya@bima.com")
	req.Header.Add("X-User-Role", "2")

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
