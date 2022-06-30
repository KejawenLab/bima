package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/utils"
	"github.com/golang-jwt/jwt/v4"
)

type Jwt struct {
	Secret string
	Method string
	User   *configs.User
}

func (j *Jwt) Attach(request *http.Request, response http.ResponseWriter) bool {
	ctx := context.WithValue(context.Background(), "scope", "auth_middleware")
	bearerToken := strings.Split(strings.TrimSpace(request.Header.Get("Authorization")), " ")
	if len(bearerToken) != 2 {
		loggers.Logger.Error(ctx, "token not provided")
		http.Error(response, "unauthorization", http.StatusUnauthorized)

		return true
	}

	token, err := utils.ValidateToken(j.Secret, j.Method, bearerToken[1])
	if err != nil {
		loggers.Logger.Error(ctx, err.Error())
		http.Error(response, "unauthorization", http.StatusUnauthorized)

		return true
	}

	claims, _ := token.(jwt.MapClaims)
	if id, ok := claims["id"]; ok {
		j.User.Id = id.(string)
	}

	if email, ok := claims["email"]; ok {
		j.User.Email = email.(string)
	}

	if role, ok := claims["role"]; ok {
		j.User.Role = int(role.(float64))
	}

	return false
}

func (j *Jwt) Priority() int {
	return 257
}
