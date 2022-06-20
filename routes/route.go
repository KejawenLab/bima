package routes

import (
	"net/http"

	"github.com/KejawenLab/bima/v2/middlewares"
	"google.golang.org/grpc"
)

type Route interface {
	Path() string
	Method() string
	Handle(w http.ResponseWriter, r *http.Request, params map[string]string)
	SetClient(client *grpc.ClientConn)
	Middlewares() []middlewares.Middleware
}
