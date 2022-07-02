package utils

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func Test_Jwt(t *testing.T) {
	_, err := ValidateToken("secret", jwt.SigningMethodHS512.Name, "invalid")
	assert.NotNil(t, err)

	claims := jwt.MapClaims{
		"id":    "test",
		"email": "test@mail.com",
		"role":  1,
	}

	token, err := CreateToken("secret", jwt.SigningMethodHS512.Name, claims, 2)

	assert.Nil(t, err)
	assert.NotNil(t, token)

	_, err = ValidateToken("secret", "invalid", token)
	assert.NotNil(t, err)

	result, err := ValidateToken("secret", jwt.SigningMethodHS512.Name, token)
	assert.Nil(t, err)

	assert.Equal(t, claims["id"], result["id"])
	assert.Equal(t, claims["email"], result["email"])
	assert.Equal(t, claims["role"], int(result["role"].(float64)))
}

func Test_Refresh_Jwt(t *testing.T) {
	_, err := ValidateRefreshToken("secret", jwt.SigningMethodHS512.Name, "invalid")
	assert.NotNil(t, err)

	claims := jwt.MapClaims{
		"id":    "test",
		"email": "test@mail.com",
		"role":  1,
	}

	token, err := CreateToken("secret", jwt.SigningMethodHS512.Name, claims, 2)

	refreshToken, err := CreateRefreshToken("secret", jwt.SigningMethodHS512.Name, token)

	assert.Nil(t, err)
	assert.NotNil(t, token)

	_, err = ValidateRefreshToken("secret", "invalid", refreshToken)
	assert.NotNil(t, err)

	_, err = ValidateRefreshToken("secret", jwt.SigningMethodHS512.Name, token)
	assert.NotNil(t, err)

	result, err := ValidateRefreshToken("secret", jwt.SigningMethodHS512.Name, refreshToken)
	assert.Nil(t, err)

	assert.Equal(t, claims["id"], result["id"])
	assert.Equal(t, claims["email"], result["email"])
	assert.Equal(t, claims["role"], int(result["role"].(float64)))
}
