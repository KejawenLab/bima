package routes

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/crowdeco/bima/configs"
	"google.golang.org/grpc"
)

const API_DOC_PATH = "/api/docs/"

type ApiDoc struct {
	Env *configs.Env
}

func (a *ApiDoc) Path() string {
	return fmt.Sprintf("%s{path}", API_DOC_PATH)
}

func (a *ApiDoc) Method() string {
	return http.MethodGet
}

func (a *ApiDoc) SetClient(client *grpc.ClientConn) {
}

func (a *ApiDoc) Handle(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if a.Env.Debug {
		regex := regexp.MustCompile(API_DOC_PATH)
		http.ServeFile(w, r, regex.ReplaceAllString(r.URL.Path, "swaggers/"))
	}
}
