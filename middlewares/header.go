package middlewares

import (
	"net/http"

	"github.com/KejawenLab/bima/v2"
	configs "github.com/KejawenLab/bima/v2/configs"
)

type Header struct {
}

func (a *Header) Attach(_ *http.Request, response http.ResponseWriter) bool {
	response.Header().Add("X-Bima-Version", bima.VERSION_STRING)

	return false
}

func (a *Header) Priority() int {
	return configs.LOWEST_PRIORITY - 1
}
