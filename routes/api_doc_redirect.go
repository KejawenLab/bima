package routes

import (
	"fmt"
	"net/http"

	"github.com/KejawenLab/bima/v2/middlewares"
	"google.golang.org/grpc"
)

type ApiDocRedirect struct {
}

func (a *ApiDocRedirect) Path() string {
	return API_DOC_PATH
}

func (a *ApiDocRedirect) Method() string {
	return http.MethodGet
}

func (a *ApiDocRedirect) SetClient(client *grpc.ClientConn) {}

func (a *ApiDocRedirect) Middlewares() []middlewares.Middleware {
	return nil
}

func (a *ApiDocRedirect) Handle(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	http.Redirect(w, r, fmt.Sprintf("%s/", r.URL.RequestURI()), http.StatusPermanentRedirect)
}
