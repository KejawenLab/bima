package handlers

import (
	"net/http"
	"sort"

	configs "github.com/crowdeco/bima/configs"
)

type Middleware struct {
	Middlewares []configs.Middleware
}

func (m *Middleware) Register(middlewares []configs.Middleware) {
	sort.Slice(middlewares, func(i, j int) bool {
		return middlewares[i].Priority() > middlewares[j].Priority()
	})
	m.Middlewares = middlewares
}

func (m *Middleware) Attach(handler http.Handler) http.Handler {
	sort.Slice(m.Middlewares, func(i, j int) bool {
		return m.Middlewares[i].Priority() > m.Middlewares[j].Priority()
	})

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		for _, middleware := range m.Middlewares {
			stop := middleware.Attach(request, response)
			if stop {
				return
			}
		}

		handler.ServeHTTP(response, request)
	})
}
