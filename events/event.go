package events

import (
	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/paginations"
	"github.com/kamva/mgm/v3"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

const (
	PaginationEvent   = Event("pagination")
	BeforeCreateEvent = Event("before_create")
	BeforeUpdateEvent = Event("before_update")
	BeforeDeleteEvent = Event("before_delete")
	AfterCreateEvent  = Event("after_create")
	AfterUpdateEvent  = Event("after_update")
	AfterDeleteEvent  = Event("after_delete")
)

type (
	Event string

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

func (e Event) String() string {
	return string(e)
}
