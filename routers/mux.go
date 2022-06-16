package routers

import (
	"context"
	"net/http"

	"github.com/KejawenLab/bima/v2"
	"github.com/KejawenLab/bima/v2/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type MuxRouter struct {
	routes []configs.Route
}

func (m *MuxRouter) Register(routes []configs.Route) {
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
	return bima.LowestPriority - 1
}
