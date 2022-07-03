package middlewares

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func Test_Jwt(t *testing.T) {
	loggers.Default("test")

	config := &configs.Env{}
	middleware := Jwt{
		Debug:         true,
		Secret:        "secret",
		SigningMethod: jwt.SigningMethodHS512.Name,
		Env:           config,
		Whitelist:     "/bar",
	}

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, true, middleware.Attach(req, w))

	req = httptest.NewRequest("GET", "http://bima.framework/bar", nil)
	w = httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))

	req = httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	req.Header.Add("Authorization", "Bearer invalid")
	w = httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, true, middleware.Attach(req, w))

	claims := jwt.MapClaims{
		"id": "test",
	}

	token, err := utils.CreateToken(middleware.Secret, middleware.SigningMethod, claims, 2)

	assert.Nil(t, err)

	req = httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, true, middleware.Attach(req, w))
	assert.Empty(t, req.Header.Get("user"))

	claims = jwt.MapClaims{
		"user": "test",
	}

	token, err = utils.CreateToken(middleware.Secret, middleware.SigningMethod, claims, 2)

	assert.Nil(t, err)

	req = httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
	assert.NotEmpty(t, req.Header.Get("X-Bima-User"))
	assert.Equal(t, req.Header.Get("X-Bima-User"), config.User)
}
