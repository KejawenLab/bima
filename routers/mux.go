package routers

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/KejawenLab/bima/v4/loggers"
	"github.com/KejawenLab/bima/v4/routes"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type MuxRouter struct {
	Debug  bool
	routes []routes.Route
}

func (m *MuxRouter) Register(muxs []routes.Route) {
	for _, mux := range muxs {
		m.Add(mux)
	}
}

func (m *MuxRouter) Add(route routes.Route) {
	m.routes = append(m.routes, route)
}

func (m *MuxRouter) Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) {
	for _, v := range m.routes {
		route := v
		route.SetClient(client)
		server.HandlePath(route.Method(), route.Path(), func(w http.ResponseWriter, r *http.Request, params map[string]string) {
			if !m.Debug {
				for _, middleware := range route.Middlewares() {
					if stop := middleware.Attach(r, w); stop {
						return
					}
				}

				route.Handle(w, r, params)

				return
			}

			for _, middleware := range route.Middlewares() {
				if stop := middleware.Attach(r, w); stop {
					var stopper strings.Builder
					stopper.WriteString("middleware stopped by: ")
					stopper.WriteString(reflect.TypeOf(middleware).Elem().Name())

					loggers.Logger.Debug(context, stopper.String())

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
