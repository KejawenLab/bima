package middlewares

import (
	"net/http"

	"github.com/rs/cors"
)

type Cors struct {
	Options cors.Options
}

func (c *Cors) Attach(request *http.Request, response http.ResponseWriter) bool {
	handler := cors.New(c.Options)
	handler.HandlerFunc(response, request)

	return false
}

func (c *Cors) Priority() int {
	return -255
}
