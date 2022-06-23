package interfaces

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/middlewares"
	"github.com/KejawenLab/bima/v3/routers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/grpclog"
)

type Rest struct {
	GRpcPort   int
	HttpPort   int
	Middleware *middlewares.Factory
	Router     *routers.Factory
}

func (r *Rest) Run(servers []configs.Server) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	endpoint := fmt.Sprintf("0.0.0.0:%d", r.GRpcPort)
	gRpcClient, err := grpc.DialContext(ctx, endpoint, options...)
	if err != nil {
		log.Fatalf("Server is not ready. %v", err)
	}

	go func() {
		<-ctx.Done()
		if cerr := gRpcClient.Close(); cerr != nil {
			grpclog.Infof("Error closing connection to %s: %v", endpoint, cerr)
		}
	}()

	r.Middleware.Sort()
	r.Router.Sort()

	http.ListenAndServe(fmt.Sprintf(":%d", r.HttpPort), r.Middleware.Attach(r.Router.Handle(ctx, http.NewServeMux(), gRpcClient)))
}

func (r *Rest) IsBackground() bool {
	return false
}

func (r *Rest) Priority() int {
	return -253
}
