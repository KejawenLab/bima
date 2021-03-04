package middlewares

import (
	"fmt"
	"net/http"

	configs "github.com/crowdeco/bima/configs"
	"github.com/crowdeco/bima/handlers"
)

type Recovery struct {
	Logger *handlers.Logger
}

func (r *Recovery) Attach(request *http.Request, response http.ResponseWriter) bool {
	defer func() {
		rc := recover()
		if rc != nil {
			switch x := rc.(type) {
			case string:
				r.Logger.Error(x)
			case error:
				r.Logger.Error(x.Error())
			default:
				r.Logger.Error(fmt.Sprintf("%+v\n", rc))
			}
		}
	}()

	return false
}

func (r *Recovery) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
