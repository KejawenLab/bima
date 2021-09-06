package routes

import (
	"context"

	configs "github.com/crowdeco/bima/v2/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GRpcGateway struct {
	Servers []configs.Server
}

func (g *GRpcGateway) Register(servers []configs.Server) {
	g.Servers = servers
}

func (g *GRpcGateway) Handle(ctx context.Context, server *runtime.ServeMux, client *grpc.ClientConn) {
	for _, handler := range g.Servers {
		handler.GRpcHandler(ctx, server, client)
	}
}

func (a *GRpcGateway) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
