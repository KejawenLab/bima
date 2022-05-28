package middlewares

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"time"

	configs "github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
)

type RequestID struct {
	Logger *handlers.Logger
}

func (r *RequestID) Attach(request *http.Request, response http.ResponseWriter) bool {
	reqID := fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().Format(time.RFC3339))))

	request.Header.Set("X-Request-Id", reqID)
	r.Logger.RequestID = reqID

	return false
}

func (r *RequestID) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
