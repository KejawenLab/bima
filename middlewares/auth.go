package middlewares

import (
	"net/http"
	"regexp"
	"strconv"

	configs "github.com/KejawenLab/bima/v2/configs"
)

type Auth struct {
	Env *configs.Env
}

func (a *Auth) Attach(request *http.Request, response http.ResponseWriter) bool {
	match, _ := regexp.MatchString(a.Env.AuthHeader.Whitelist, request.RequestURI)
	if match {
		return false
	}

	a.Env.User.Id = request.Header.Get(a.Env.AuthHeader.Id)
	a.Env.User.Email = request.Header.Get(a.Env.AuthHeader.Email)
	a.Env.User.Role, _ = strconv.Atoi(request.Header.Get(a.Env.AuthHeader.Role))

	if a.Env.User.Role == 0 || a.Env.User.Role > a.Env.AuthHeader.MaxRole {
		http.Error(response, "Unauthorization", http.StatusUnauthorized)

		return true
	}

	return false
}

func (a *Auth) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
