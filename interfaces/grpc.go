package interfaces

import (
	"bytes"
	"log"
	"net"
	"strconv"

	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/loggers"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GRpc struct {
	GRpcPort int
	Debug    bool
}

func (g *GRpc) Run(servers []configs.Server) {
	var gRpcAddress bytes.Buffer
	gRpcAddress.WriteString(":")
	gRpcAddress.WriteString(strconv.Itoa(g.GRpcPort))

	listen, err := net.Listen("tcp", gRpcAddress.String())
	if err != nil {
		log.Fatalf("Port %d is not available. %v", g.GRpcPort, err)
	}

	streams := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(),
	}
	unaries := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
	}
	if g.Debug {
		options := []grpc_logrus.Option{
			grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
		}
		streams = append(streams, grpc_logrus.StreamServerInterceptor(logrus.NewEntry(loggers.Logger.Engine), options...))
		unaries = append(unaries, grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(loggers.Logger.Engine), options...))
	}

	gRpc := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streams...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaries...)),
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
