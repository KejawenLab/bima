package middlewares

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"time"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
)

type RequestID struct {
	Logger *handlers.Logger
	Env    *configs.Env
}

func (r *RequestID) Attach(request *http.Request, response http.ResponseWriter) bool {
	reqID := request.Header.Get(r.Env.RequestIDHeader)
	if reqID == "" {
		reqID = fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().Format(time.RFC3339))))
	}

	response.Header().Add(r.Env.RequestIDHeader, reqID)
	r.Logger.RequestID = reqID

	return false
}

func (r *RequestID) Priority() int {
	return configs.LOWEST_PRIORITY - 1
}
