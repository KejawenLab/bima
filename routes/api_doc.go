package routes

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/KejawenLab/bima/v2/middlewares"
	"google.golang.org/grpc"
)

const API_DOC_PATH = "/api/docs"

type ApiDoc struct {
	Debug bool
}

func (a *ApiDoc) Path() string {
	return fmt.Sprintf("%s/{path}", API_DOC_PATH)
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

	regex := regexp.MustCompile(fmt.Sprintf("%s/", API_DOC_PATH))
	http.ServeFile(w, r, regex.ReplaceAllString(r.URL.Path, "swaggers/"))
}
