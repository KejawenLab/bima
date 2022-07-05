package middlewares

import (
	"net/http"

	"github.com/KejawenLab/bima/v3"
)

type Header struct {
}

func (h *Header) Attach(_ *http.Request, response http.ResponseWriter) bool {
	response.Header().Add("X-Bima-Version", bima.Version)

	return false
}

func (h *Header) Priority() int {
	return -257
}
