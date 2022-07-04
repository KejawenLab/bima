package middlewares

import (
	"context"
	"net/http"

	"github.com/KejawenLab/bima/v3/loggers"
)

type (
	BasicAuth struct {
		Validator ValidateUsernameAndPassword
	}

	ValidateUsernameAndPassword func(username string, password string) bool
)

func (b *BasicAuth) Attach(request *http.Request, response http.ResponseWriter) bool {
	ctx := context.WithValue(context.Background(), "scope", "basic_auth_middleware")
	username, password, ok := request.BasicAuth()
	if !ok {
		loggers.Logger.Error(ctx, "error parsing basic auth")
		http.Error(response, "invalid username or password", http.StatusUnauthorized)

		return true
	}

	if !b.Validator(username, password) {
		loggers.Logger.Error(ctx, "invalid username or password")
		http.Error(response, "invalid username or password", http.StatusUnauthorized)

		return true
	}

	return false
}

func (b *BasicAuth) Priority() int {
	return 257
}
