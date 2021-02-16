package middlewares

import (
	"net/http"

	"github.com/crowdeco/bima"
	configs "github.com/crowdeco/bima/configs"
)

type Version struct {
}

func (v *Version) Attach(request *http.Request, response http.ResponseWriter) bool {
	response.Header().Add("X-Bima-Version", bima.VERSION_STRING)

	return false
}

func (v *Version) Priority() int {
	return configs.LOWEST_PRIORITY - 1
}
