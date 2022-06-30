package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateToken(secret string, signing string, bearerToken string) (jwt.Claims, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method != jwt.GetSigningMethod(signing) {
			return nil, errors.New("invalid token")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims, err
}

func CreateToken(secret string, method string, claims jwt.MapClaims, expireAt int) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(method))

	claims["exp"] = time.Now().Add(time.Hour * time.Duration(expireAt)).Unix()
	token.Claims = claims

	return token.SignedString([]byte(secret))
}
