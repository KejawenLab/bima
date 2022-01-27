package routes

import (
	"encoding/json"
	"net/http"

	"github.com/KejawenLab/bima/v2"
	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
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
		"name":    "Bima",
		"author":  "Muhamad Surya Iksanudin<surya.iksanudin@gmail.com>",
		"link":    "https://github.com/KejawenLab/skeleton",
	}

	json.NewEncoder(w).Encode(payload)
}
