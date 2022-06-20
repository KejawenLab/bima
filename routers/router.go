package routers

import (
	"context"
	"net/http"
	"sort"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type (
	Router interface {
		Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn)
		Priority() int
	}

	Factory struct {
		Routers []Router
	}
)

func (r *Factory) Sort() {
	sort.Slice(r.Routers, func(i int, j int) bool {
		return r.Routers[i].Priority() > r.Routers[j].Priority()
	})
}

func (r *Factory) Handle(context context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux {
	mux := runtime.NewServeMux()
	for _, route := range r.Routers {
		route.Handle(context, mux, client)
	}

	server.Handle("/", mux)

	return server
}
