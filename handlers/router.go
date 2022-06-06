package handlers

import (
	"context"
	"net/http"
	"sort"

	configs "github.com/KejawenLab/bima/v2/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Router struct {
	Routers []configs.Router
}

func (r *Router) Sort() {
	sort.Slice(r.Routers, func(i int, j int) bool {
		return r.Routers[i].Priority() > r.Routers[j].Priority()
	})
}

func (r *Router) Handle(context context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux {
	mux := runtime.NewServeMux()
	for _, route := range r.Routers {
		route.Handle(context, mux, client)
	}

	server.Handle("/", mux)

	return server
}
