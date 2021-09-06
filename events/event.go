package events

import (
	"net/http"

	paginations "github.com/KejawenLab/bima/v2/paginations"
	services "github.com/KejawenLab/bima/v2/services"
	elastic "github.com/olivere/elastic/v7"
)

const PAGINATION_EVENT = "event.pagination"
const BEFORE_CREATE_EVENT = "event.before_create"
const AFTER_CREATE_EVENT = "event.after_create"
const BEFORE_UPDATE_EVENT = "event.before_update"
const AFTER_UPDATE_EVENT = "event.after_update"
const BEFORE_DELETE_EVENT = "event.before_delete"
const AFTER_DELETE_EVENT = "event.after_delete"
const REQUEST_EVENT = "event.request"
const RESPONSE_EVENT = "event.response"

type (
	Model struct {
		Data       interface{}
		Id         string
		Repository *services.Repository
	}

	Pagination struct {
		Repository *services.Repository
		Query      *elastic.BoolQuery
		Filters    []paginations.Filter
	}

	Request struct {
		HttpRequest *http.Request
	}

	Response struct {
		ResponseWriter http.ResponseWriter
	}
)
