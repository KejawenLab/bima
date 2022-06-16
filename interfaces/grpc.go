package interfaces

import (
	"fmt"
	"log"
	"net"

	"github.com/KejawenLab/bima/v2"
	"github.com/KejawenLab/bima/v2/configs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

type GRpc struct {
	GRpcPort int
}

func (g *GRpc) Run(servers []configs.Server) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", g.GRpcPort))
	if err != nil {
		log.Fatalf("Port %d is not available. %v", g.GRpcPort, err)
	}

	gRpc := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	for _, server := range servers {
		server.RegisterGRpc(gRpc)
	}

	gRpc.Serve(listen)
}

func (g *GRpc) IsBackground() bool {
	return true
}

func (g *GRpc) Priority() int {
	return bima.HighestPriority + 1
}
