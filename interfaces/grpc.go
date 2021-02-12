package interfaces

import (
	"fmt"
	"log"
	"net"

	configs "github.com/crowdeco/bima/configs"
	"github.com/fatih/color"
	grpc "google.golang.org/grpc"
)

type GRpc struct {
	Env  *configs.Env
	GRpc *grpc.Server
}

func (g *GRpc) Run(servers []configs.Server) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Env.RpcPort))
	if err != nil {
		log.Fatalf("Port %d is not available. %v", g.Env.RpcPort, err)
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
