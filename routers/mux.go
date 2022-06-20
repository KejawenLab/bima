package routers

import (
	"context"
	"net/http"

	"github.com/KejawenLab/bima/v2/routes"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type MuxRouter struct {
	routes []routes.Route
}

func (m *MuxRouter) Register(routes []routes.Route) {
	m.routes = append(m.routes, routes...)
}

func (m *MuxRouter) Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) {
	for _, v := range m.routes {
		route := v
		route.SetClient(client)
		server.HandlePath(route.Method(), route.Path(), func(w http.ResponseWriter, r *http.Request, params map[string]string) {
			for _, m := range v.Middlewares() {
				if stop := m.Attach(r, w); stop {
					return
				}
			}

			route.Handle(w, r, params)
		})
	}
}

func (m *MuxRouter) Priority() int {
	return -255
}
