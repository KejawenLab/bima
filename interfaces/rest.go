package interfaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type Rest struct {
	Env        *configs.Env
	Middleware *handlers.Middleware
	Router     *handlers.Router
	Server     *http.ServeMux
	Context    context.Context
}

func (r *Rest) Run(servers []configs.Server) {
	ctx, cancel := context.WithCancel(r.Context)
	defer cancel()

	endpoint := fmt.Sprintf("0.0.0.0:%d", r.Env.RpcPort)
	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close connection to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close connection to %s: %v", endpoint, cerr)
			}
		}()
	}()

	util := color.New(color.FgCyan, color.Bold)

	util.Printf("✓ ")
	fmt.Printf("Playing REST Multimedia on :%d...\n", r.Env.HtppPort)

	util.Printf("✓ ")
	fmt.Println("Playlist API is Ready at /api/docs...")

	r.Middleware.Sort()
	r.Router.Sort()

	http.ListenAndServe(fmt.Sprintf(":%d", r.Env.HtppPort), r.Middleware.Attach(r.Router.Handle(ctx, r.Server, conn)))
}

func (r *Rest) IsBackground() bool {
	return false
}

func (r *Rest) Priority() int {
	return configs.LOWEST_PRIORITY - 1
}
