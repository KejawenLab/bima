package routers

import (
	"context"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GRpcGateway struct {
	servers []configs.Server
}

func (g *GRpcGateway) Register(servers []configs.Server) {
	g.servers = servers
}

func (g *GRpcGateway) Handle(ctx context.Context, server *runtime.ServeMux, client *grpc.ClientConn) {
	for _, handler := range g.servers {
		handler.GRpcHandler(ctx, server, client)
	}
}

func (a *GRpcGateway) Priority() int {
	return 255
}
