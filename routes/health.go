package routes

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const HEALTH_PATH = "/health"

type Health struct {
	Client *grpc.ClientConn
}

func (h *Health) Path() string {
	return HEALTH_PATH
}

func (h *Health) Method() string {
	return http.MethodGet
}

func (h *Health) SetClient(client *grpc.ClientConn) {
	h.Client = client
}

func (h *Health) Handle(w http.ResponseWriter, r *http.Request, params map[string]string) {
	w.Header().Set("Content-Type", "text/html")
	s := h.Client.GetState()

	if s != connectivity.Ready {
		http.Error(w, fmt.Sprintf("gRPC server is %s", s), http.StatusBadGateway)

		return
	}

	fmt.Fprintln(w, "OK")
}
