package interfaces

import (
	"bytes"
	"log"
	"net"
	"strconv"

	"github.com/KejawenLab/bima/v3/configs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

type GRpc struct {
	GRpcPort int
}

func (g *GRpc) Run(servers []configs.Server) {
	var gRpcAddress bytes.Buffer
	gRpcAddress.WriteString(":")
	gRpcAddress.WriteString(strconv.Itoa(g.GRpcPort))

	listen, err := net.Listen("tcp", gRpcAddress.String())
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
	return 257
}
