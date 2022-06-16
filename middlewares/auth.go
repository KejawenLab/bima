package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
)

type Auth struct {
	Env    *configs.Env
	Logger *handlers.Logger
}

func (a *Auth) Attach(request *http.Request, response http.ResponseWriter) bool {
	ctx := context.WithValue(context.Background(), "scope", "auth_middleware")
	if a.Env.AuthHeader.Id == "" || a.Env.AuthHeader.Email == "" || a.Env.AuthHeader.Role == "" {
		a.Logger.Debug(ctx, "Auth header not set in environment variables")

		return false
	}

	match, _ := regexp.MatchString(a.Env.AuthHeader.Whitelist, request.RequestURI)
	if match {
		a.Logger.Debug(ctx, fmt.Sprintf("Whitelisting url %s", request.RequestURI))

		return false
	}

	a.Env.User = &configs.User{}
	a.Env.User.Id = request.Header.Get(a.Env.AuthHeader.Id)
	a.Env.User.Email = request.Header.Get(a.Env.AuthHeader.Email)
	a.Env.User.Role, _ = strconv.Atoi(request.Header.Get(a.Env.AuthHeader.Role))
	if a.Env.User.Id == "" || a.Env.User.Email == "" || a.Env.User.Role == 0 {
		a.Logger.Error(ctx, "User is not provided")
		http.Error(response, "Unauthorization", http.StatusUnauthorized)

		return true
	}

	if a.Env.User.Role < a.Env.AuthHeader.MinRole {
		a.Logger.Error(ctx, "Insufficient role")
		http.Error(response, "Unauthorization", http.StatusUnauthorized)

		return true
	}

	return false
}

func (a *Auth) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
