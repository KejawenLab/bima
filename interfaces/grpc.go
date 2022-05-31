package interfaces

import (
	"fmt"
	"net"

	configs "github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/fatih/color"
	grpc "google.golang.org/grpc"
)

type GRpc struct {
	Env    *configs.Env
	GRpc   *grpc.Server
	Logger *handlers.Logger
}

func (g *GRpc) Run(servers []configs.Server) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Env.RpcPort))
	if err != nil {
		g.Logger.Fatal(fmt.Sprintf("Port %d is not available. %v", g.Env.RpcPort, err))
	}

	for _, server := range servers {
		server.RegisterGRpc(g.GRpc)
	}

	color.New(color.FgCyan, color.Bold).Printf("✓ ")
	fmt.Printf("Connecting gRPC Multimedia on :%d...\n", g.Env.RpcPort)

	g.GRpc.Serve(l)
}

func (g *GRpc) IsBackground() bool {
	return true
}

func (g *GRpc) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
