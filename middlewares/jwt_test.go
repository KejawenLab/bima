package middlewares

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func Test_Jwt(t *testing.T) {
	middleware := Jwt{
		Secret: "secret",
		Method: jwt.SigningMethodHS512.Name,
		User:   &configs.User{},
	}

	req := httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, true, middleware.Attach(req, w))

	req = httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	req.Header.Add("Authorization", "Bearer invalid")
	w = httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, true, middleware.Attach(req, w))

	claims := jwt.MapClaims{
		"id":    "test",
		"email": "test@mail.com",
		"role":  1,
	}

	token, err := utils.CreateToken(middleware.Secret, middleware.Method, claims, 2)

	assert.Nil(t, err)

	req = httptest.NewRequest("GET", "http://bima.framework/foo", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w = httptest.NewRecorder()

	assert.Equal(t, 257, middleware.Priority())
	assert.Equal(t, false, middleware.Attach(req, w))
}
