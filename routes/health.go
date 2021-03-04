package routes

import (
	"encoding/json"
	"net/http"

	"github.com/crowdeco/bima"
	"github.com/crowdeco/bima/configs"
	"github.com/crowdeco/bima/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const HEALTH_PATH = "/health"

type Health struct {
	Client *grpc.ClientConn
	Logger *handlers.Logger
}

func (h *Health) Path() string {
	return HEALTH_PATH
}

func (h *Health) Method() string {
	return http.MethodGet
}

func (h *Health) Middlewares() []configs.Middleware {
	return nil
}

func (h *Health) SetClient(client *grpc.ClientConn) {
	h.Client = client
}

func (h *Health) Handle(w http.ResponseWriter, r *http.Request, params map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	s := h.Client.GetState()

	if s != connectivity.Ready {
		h.Logger.Error("gRPC server is down")
		http.Error(w, "gRPC server is down", http.StatusBadGateway)

		return
	}

	payload := map[string]string{
		"version": bima.VERSION_STRING,
		"x_x":     "01001011 01100101 01101010 01100001 01110111 01100101 01101110 01001100 01100001 01100010",
		"v_v":     "01001101 01110101 01101000 01100001 01101101 01100001 01100100 00100000 01010011 01110101 01110010 01111001 01100001 00100000 01001001 01101011 01110011 01100001 01101110 01110101 01100100 01101001 01101110",
	}

	json.NewEncoder(w).Encode(payload)
}
