package events

import "net/http"

type Request struct {
	HttpRequest *http.Request
}
