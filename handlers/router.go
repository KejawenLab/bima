package handlers

import (
	"context"
	"net/http"
	"sort"

	configs "github.com/Kejawenlab/bima/v2/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Router struct {
	Routes []configs.Router
}

func (r *Router) Handle(context context.Context, server *http.ServeMux, client *grpc.ClientConn) *http.ServeMux {
	sort.Slice(r.Routes, func(i int, j int) bool {
		return r.Routes[i].Priority() > r.Routes[j].Priority()
	})

	mux := runtime.NewServeMux()
	for _, route := range r.Routes {
		route.Handle(context, mux, client)
	}

	server.Handle("/", mux)

	return server
}
