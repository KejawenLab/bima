package middlewares

import (
	"net/http"

	"github.com/KejawenLab/bima/v2"
	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/google/uuid"
)

type RequestID struct {
	Logger          *handlers.Logger
	RequestIDHeader string
}

func (r *RequestID) Attach(request *http.Request, response http.ResponseWriter) bool {
	reqID := request.Header.Get(r.RequestIDHeader)
	if reqID == "" {
		reqID = uuid.NewString()
	}

	response.Header().Add(r.RequestIDHeader, reqID)
	r.Logger.Add("request_id", reqID)

	return false
}

func (r *RequestID) Priority() int {
	return bima.HighestPriority + 2
}
