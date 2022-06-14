package interfaces

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type Rest struct {
	GRpcPort   int
	HttpPort   int
	Middleware *handlers.Middleware
	Router     *handlers.Router
}

func (r *Rest) Run(servers []configs.Server) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := fmt.Sprintf("0.0.0.0:%d", r.GRpcPort)
	gRpcClient, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Server is not ready. %v", err)
	}

	defer func() {
		if err != nil {
			if cerr := gRpcClient.Close(); cerr != nil {
				grpclog.Infof("Failed to close connection to %s: %v", endpoint, cerr)
			}
			return
		}

		go func() {
			<-ctx.Done()
			if cerr := gRpcClient.Close(); cerr != nil {
				grpclog.Infof("Context closed by %s: %v", endpoint, cerr)
			}
		}()
	}()

	r.Middleware.Sort()
	r.Router.Sort()

	http.ListenAndServe(fmt.Sprintf(":%d", r.HttpPort), r.Middleware.Attach(r.Router.Handle(ctx, http.NewServeMux(), gRpcClient)))
}

func (r *Rest) IsBackground() bool {
	return false
}

func (r *Rest) Priority() int {
	return configs.LOWEST_PRIORITY - 1
}
