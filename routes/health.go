package routes

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/goccy/go-json"

	"github.com/KejawenLab/bima/v2"
	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const HelthPath = "/health"

type Health struct {
	client *grpc.ClientConn
	Logger *handlers.Logger
}

func (h *Health) Path() string {
	return HelthPath
}

func (h *Health) Method() string {
	return http.MethodGet
}

func (h *Health) Middlewares() []configs.Middleware {
	return nil
}

func (h *Health) SetClient(client *grpc.ClientConn) {
	h.client = client
}

func (h *Health) Handle(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	s := h.client.GetState()
	if s != connectivity.Ready {
		h.Logger.Error(context.WithValue(context.Background(), "scope", "health_route"), "gRPC server is down")
		http.Error(w, "gRPC server is down", http.StatusBadGateway)

		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	payload := map[string]interface{}{
		"version": bima.Version,
		"name":    "Bima",
		"author":  "Muhamad Surya Iksanudin<surya.iksanudin@gmail.com>",
		"link":    "https://github.com/KejawenLab/skeleton",
		"memory_usage": map[string]string{
			"allocation":       fmt.Sprintf("%d MiB", h.byteToMb(m.Alloc)),
			"total_allocation": fmt.Sprintf("%d MiB", h.byteToMb(m.TotalAlloc)),
			"system":           fmt.Sprintf("%d MiB", h.byteToMb(m.Sys)),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func (h *Health) byteToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
