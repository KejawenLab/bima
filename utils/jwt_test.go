package utils

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func Test_Jwt(t *testing.T) {
	result, err := ValidateToken("secret", jwt.SigningMethodHS512.Name, "invalid")
	assert.NotNil(t, err)
	assert.Nil(t, result)

	claims := jwt.MapClaims{
		"id":    "test",
		"email": "test@mail.com",
		"role":  1,
	}

	token, err := CreateToken("secret", jwt.SigningMethodHS512.Name, claims, 2)

	assert.Nil(t, err)
	assert.NotNil(t, token)

	result, err = ValidateToken("secret", "invalid", token)
	assert.NotNil(t, err)
	assert.Nil(t, result)

	result, err = ValidateToken("secret", jwt.SigningMethodHS512.Name, token)
	assert.Nil(t, err)

	mapClaims, ok := result.(jwt.MapClaims)

	assert.True(t, ok)
	assert.Equal(t, claims["id"], mapClaims["id"])
	assert.Equal(t, claims["email"], mapClaims["email"])
	assert.Equal(t, claims["role"], int(mapClaims["role"].(float64)))
}
