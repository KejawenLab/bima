package events

import (
	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/paginations"
	"github.com/kamva/mgm/v3"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

const PAGINATION_EVENT = "event.pagination"
const BEFORE_CREATE_EVENT = "event.before_create"
const AFTER_CREATE_EVENT = "event.after_create"
const BEFORE_UPDATE_EVENT = "event.before_update"
const AFTER_UPDATE_EVENT = "event.after_update"
const BEFORE_DELETE_EVENT = "event.before_delete"
const AFTER_DELETE_EVENT = "event.after_delete"

type (
	Model struct {
		Data       interface{}
		Id         string
		Repository configs.Repository
	}

	ElasticsearchPagination struct {
		Query   *elastic.BoolQuery
		Filters []paginations.Filter
	}

	MongodbPagination struct {
		Query         *mgm.Collection
		Filters       []paginations.Filter
		MongoDbFilter bson.M
	}

	GormPagination struct {
		Query   *gorm.DB
		Filters []paginations.Filter
	}
)
