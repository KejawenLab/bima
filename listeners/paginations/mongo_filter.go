package paginations

import (
	"strings"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoDbFilter struct {
}

func (u *MongoDbFilter) Handle(event interface{}) interface{} {
	e, ok := event.(*events.MongodbPagination)
	if !ok {
		return event
	}

	bFilters := bson.M{}
	for _, f := range e.Filters {
		bFilters[strings.ToLower(f.Field)] = bson.M{
			operator.Regex: primitive.Regex{
				Pattern: f.Value,
				Options: "im",
			},
		}
	}

	e.MongoDbFilter = bFilters

	return e
}

func (u *MongoDbFilter) Listen() string {
	return events.PAGINATION_EVENT
}

func (u *MongoDbFilter) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
