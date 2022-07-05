package routes

import (
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/goccy/go-json"

	"github.com/KejawenLab/bima/v4"
	"github.com/KejawenLab/bima/v4/middlewares"
	"google.golang.org/grpc"
)

const HelthPath = "/health"

type Health struct {
}

func (h *Health) Path() string {
	return HelthPath
}

func (h *Health) Method() string {
	return http.MethodGet
}

func (h *Health) Middlewares() []middlewares.Middleware {
	return nil
}

func (h *Health) SetClient(client *grpc.ClientConn) {
}

func (h *Health) Handle(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var m runtime.MemStats
	var alloc strings.Builder
	var totalAlloc strings.Builder
	var system strings.Builder

	runtime.ReadMemStats(&m)

	alloc.WriteString(strconv.Itoa(h.byteToMb(m.Alloc)))
	alloc.WriteString(" MiB")
	totalAlloc.WriteString(strconv.Itoa(h.byteToMb(m.TotalAlloc)))
	totalAlloc.WriteString(" MiB")
	system.WriteString(strconv.Itoa(h.byteToMb(m.Sys)))
	system.WriteString(" MiB")

	payload := map[string]interface{}{
		"version": bima.Version,
		"name":    "Bima",
		"author":  "Muhamad Surya Iksanudin<surya.iksanudin@gmail.com>",
		"link":    "https://github.com/KejawenLab/skeleton",
		"memory_usage": map[string]string{
			"allocation":       alloc.String(),
			"total_allocation": totalAlloc.String(),
			"system":           system.String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func (h *Health) byteToMb(b uint64) int {
	return int(b / 1024 / 1024)
}
