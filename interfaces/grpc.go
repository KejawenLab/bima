package interfaces

import (
	"fmt"
	"net"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/fatih/color"
	"google.golang.org/grpc"
)

type GRpc struct {
	GRpcPort int
	GRpc     *grpc.Server
	Logger   *handlers.Logger
}

func (g *GRpc) Run(servers []configs.Server) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", g.GRpcPort))
	if err != nil {
		g.Logger.Fatal(fmt.Sprintf("Port %d is not available. %v", g.GRpcPort, err))
	}

	for _, server := range servers {
		server.RegisterGRpc(g.GRpc)
	}

	color.New(color.FgCyan, color.Bold).Printf("✓ ")
	fmt.Printf("Connecting gRPC Multimedia on :%d...\n", g.GRpcPort)

	g.GRpc.Serve(listen)
}

func (g *GRpc) IsBackground() bool {
	return true
}

func (g *GRpc) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
