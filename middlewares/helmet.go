package middlewares

import (
	"net/http"

	"github.com/goddtriffin/helmet"
)

type Helmet struct {
}

func (h *Helmet) Attach(_ *http.Request, response http.ResponseWriter) bool {
	helmet := helmet.Default()

	helmet.ContentSecurityPolicy.Header(response)
	helmet.XContentTypeOptions.Header(response)
	helmet.XDNSPrefetchControl.Header(response)
	helmet.XDownloadOptions.Header(response)
	helmet.ExpectCT.Header(response)
	helmet.FeaturePolicy.Header(response)
	helmet.XFrameOptions.Header(response)
	helmet.XPermittedCrossDomainPolicies.Header(response)
	helmet.XPoweredBy.Header(response)
	helmet.ReferrerPolicy.Header(response)
	helmet.StrictTransportSecurity.Header(response)
	helmet.XXSSProtection.Header(response)

	return false
}

func (h *Helmet) Priority() int {
	return -255
}
