package events

import (
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/repositories"
	"github.com/kamva/mgm/v3"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

type (
	Event string

	Model struct {
		Data       interface{}
		Id         string
		Repository repositories.Repository
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

const (
	PaginationEvent   = Event("pagination")
	BeforeCreateEvent = Event("before_create")
	BeforeUpdateEvent = Event("before_update")
	BeforeDeleteEvent = Event("before_delete")
	AfterCreateEvent  = Event("after_create")
	AfterUpdateEvent  = Event("after_update")
	AfterDeleteEvent  = Event("after_delete")
)

func (e Event) String() string {
	return string(e)
}
