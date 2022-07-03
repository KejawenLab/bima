package interfaces

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	var gRpcAddress strings.Builder
	gRpcAddress.WriteString("0.0.0.0:")
	gRpcAddress.WriteString(strconv.Itoa(r.GRpcPort))

	gRpcClient, err := grpc.DialContext(ctx, gRpcAddress.String(), options...)
	if err != nil {
		log.Fatalf("Server is not ready. %v", err)
	}

	go func() {
		<-ctx.Done()
		if cerr := gRpcClient.Close(); cerr != nil {
			grpclog.Infof("Error closing connection to %s: %v", gRpcAddress, cerr)
		}
	}()

	r.Middleware.Sort()
	r.Router.Sort()

	var httpAddress strings.Builder
	httpAddress.WriteString(":")
	httpAddress.WriteString(strconv.Itoa(r.HttpPort))

	http.ListenAndServe(httpAddress.String(), r.Middleware.Attach(r.Router.Handle(ctx, http.NewServeMux(), gRpcClient)))
}

func (r *Rest) IsBackground() bool {
	return false
}

func (r *Rest) Priority() int {
	return -253
}
