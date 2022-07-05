package routes

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/KejawenLab/bima/v4/middlewares"
	"google.golang.org/grpc"
)

const ApiDocPath = "/api/docs"

type ApiDoc struct {
	Debug bool
}

func (a *ApiDoc) Path() string {
	var path strings.Builder
	path.WriteString(ApiDocPath)
	path.WriteString("/{path}")

	return path.String()
}

func (a *ApiDoc) Method() string {
	return http.MethodGet
}

func (a *ApiDoc) SetClient(client *grpc.ClientConn) {}

func (a *ApiDoc) Middlewares() []middlewares.Middleware {
	return nil
}

func (a *ApiDoc) Handle(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if !a.Debug {
		w.Write([]byte("Api doc not available"))

		return
	}

	var path strings.Builder
	path.WriteString(ApiDocPath)
	path.WriteString("/")

	regex := regexp.MustCompile(path.String())
	http.ServeFile(w, r, regex.ReplaceAllString(r.URL.Path, "swaggers/"))
}
