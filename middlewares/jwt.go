package middlewares

import (
	"bytes"
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/utils"
)

type Jwt struct {
	Debug         bool
	Secret        string
	SigningMethod string
	Whitelist     string
	Env           *configs.Env
}

func (j *Jwt) Attach(request *http.Request, response http.ResponseWriter) bool {
	ctx := context.WithValue(context.Background(), "scope", "jwt_middleware")
	match, _ := regexp.MatchString(j.Whitelist, request.RequestURI)
	if match {
		if j.Debug {
			var log bytes.Buffer
			log.WriteString("whitelisting url ")
			log.WriteString(request.RequestURI)

			loggers.Logger.Debug(ctx, log.String())
		}

		return false
	}

	bearerToken := strings.Split(strings.TrimSpace(request.Header.Get("Authorization")), " ")
	if len(bearerToken) != 2 {
		loggers.Logger.Error(ctx, "token not provided")
		http.Error(response, "unauthorization", http.StatusUnauthorized)

		return true
	}

	claims, err := utils.ValidateToken(j.Secret, j.SigningMethod, bearerToken[1])
	if err != nil {
		loggers.Logger.Error(ctx, err.Error())
		http.Error(response, "unauthorization", http.StatusUnauthorized)

		return true
	}

	if user, ok := claims["user"]; ok {
		j.Env.User = user.(string)
		request.Header.Add("X-Bima-User", j.Env.User)

		return false
	}

	loggers.Logger.Error(ctx, "user not provided")
	http.Error(response, "unauthorization", http.StatusUnauthorized)

	return true
}

func (j *Jwt) Priority() int {
	return 257
}
