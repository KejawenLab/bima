package routes

import (
	"net/http"
	"strings"

	"github.com/KejawenLab/bima/v3/middlewares"
	"google.golang.org/grpc"
)

type ApiDocRedirect struct {
}

func (a *ApiDocRedirect) Path() string {
	return ApiDocPath
}

func (a *ApiDocRedirect) Method() string {
	return http.MethodGet
}

func (a *ApiDocRedirect) SetClient(client *grpc.ClientConn) {}

func (a *ApiDocRedirect) Middlewares() []middlewares.Middleware {
	return nil
}

func (a *ApiDocRedirect) Handle(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var path strings.Builder
	path.WriteString(r.URL.RequestURI())
	path.WriteString("/")

	http.Redirect(w, r, path.String(), http.StatusPermanentRedirect)
}
