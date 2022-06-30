package middlewares

import (
	"bytes"
	"context"
	"net/http"
	"regexp"
	"strconv"

	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/loggers"
)

type Auth struct {
	Env *configs.Env
}

func (a *Auth) Attach(request *http.Request, response http.ResponseWriter) bool {
	ctx := context.WithValue(context.Background(), "scope", "auth_middleware")
	if a.Env.AuthHeader.Id == "" || a.Env.AuthHeader.Email == "" || a.Env.AuthHeader.Role == "" {
		if a.Env.Debug {
			loggers.Logger.Debug(ctx, "Auth header not set in environment variables")
		}

		return false
	}

	match, _ := regexp.MatchString(a.Env.AuthHeader.Whitelist, request.RequestURI)
	if match {
		if a.Env.Debug {
			var log bytes.Buffer
			log.WriteString("Whitelisting url ")
			log.WriteString(request.RequestURI)

			loggers.Logger.Debug(ctx, log.String())
		}

		return false
	}

	a.Env.User = &configs.User{}
	a.Env.User.Id = request.Header.Get(a.Env.AuthHeader.Id)
	a.Env.User.Email = request.Header.Get(a.Env.AuthHeader.Email)
	a.Env.User.Role, _ = strconv.Atoi(request.Header.Get(a.Env.AuthHeader.Role))
	if a.Env.User.Id == "" || a.Env.User.Email == "" || a.Env.User.Role == 0 {
		loggers.Logger.Error(ctx, "User is not provided")
		http.Error(response, "Unauthorization", http.StatusUnauthorized)

		return true
	}

	if a.Env.User.Role < a.Env.AuthHeader.MinRole {
		loggers.Logger.Error(ctx, "Insufficient role")
		http.Error(response, "Unauthorization", http.StatusUnauthorized)

		return true
	}

	return false
}

func (a *Auth) Priority() int {
	return 257
}
