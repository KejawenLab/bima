package middlewares

import (
	"net/http"

	"github.com/KejawenLab/bima/v2"
)

type Header struct {
}

func (a *Header) Attach(_ *http.Request, response http.ResponseWriter) bool {
	response.Header().Add("X-Bima-Version", bima.Version)

	return false
}

func (a *Header) Priority() int {
	return bima.LowestPriority - 1
}
