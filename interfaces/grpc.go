package interfaces

import (
	"fmt"
	"log"
	"net"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/fatih/color"
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

	color.New(color.FgCyan, color.Bold).Printf("✓ ")
	fmt.Printf("Connecting gRPC Multimedia on :%d...\n", g.GRpcPort)

	gRpc.Serve(listen)
}

func (g *GRpc) IsBackground() bool {
	return true
}

func (g *GRpc) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
