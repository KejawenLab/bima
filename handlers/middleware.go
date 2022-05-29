package handlers

import (
	"net/http"
	"sort"

	configs "github.com/KejawenLab/bima/v2/configs"
	events "github.com/KejawenLab/bima/v2/events"
)

type Middleware struct {
	Dispatcher  *events.Dispatcher
	Middlewares []configs.Middleware
}

func (m *Middleware) Register(middlewares []configs.Middleware) {
	for _, v := range middlewares {
		m.Add(v)
	}
}

func (m *Middleware) Add(middlware configs.Middleware) {
	m.Middlewares = append(m.Middlewares, middlware)
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

		m.Dispatcher.Dispatch(events.REQUEST_EVENT, &events.Request{
			HttpRequest: request,
		})

		m.Dispatcher.Dispatch(events.RESPONSE_EVENT, &events.Response{
			ResponseWriter: response,
		})

		handler.ServeHTTP(response, request)
	})
}
