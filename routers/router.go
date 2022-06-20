package routers

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Router interface {
	Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn)
	Priority() int
}
