package adapter

import (
	"bytes"
	"context"

	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/kamva/mgm/v3"
	"github.com/vcraescu/go-paginator/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MongodbAdapter struct {
		Debug      bool
		Dispatcher *events.Dispatcher
	}

	mongodbPaginator struct {
		context    context.Context
		index      string
		pageQuery  *mgm.Collection
		totalQuery *mgm.Collection
		filter     bson.M
	}
)

func (mg *MongodbAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	model, ok := paginator.Model.(mgm.Model)
	if !ok {
		loggers.Logger.Fatal(ctx, "Adapter not configured properly")

		return nil
	}

	query := mgm.Coll(model)
	event := events.MongodbPagination{
		Query:         query,
		Filters:       paginator.Filters,
		MongoDbFilter: bson.M{},
	}

	if mg.Debug {
		var log bytes.Buffer
		log.WriteString("Dispatching ")
		log.WriteString(events.PaginationEvent.String())

		loggers.Logger.Debug(ctx, log.String())
	}

	mg.Dispatcher.Dispatch(events.PaginationEvent.String(), &event)

	return newMongodbPaginator(ctx, event.Query, event.MongoDbFilter)
}

func newMongodbPaginator(context context.Context, query *mgm.Collection, filter bson.M) paginator.Adapter {
	var totalQuery *mgm.Collection = query

	return &mongodbPaginator{
		context:    context,
		pageQuery:  query,
		totalQuery: totalQuery,
		filter:     filter,
	}
}

func (mg *mongodbPaginator) Nums() (int64, error) {
	return mg.totalQuery.CountDocuments(mg.context, mg.filter)
}

func (mg *mongodbPaginator) Slice(offset int, length int, data interface{}) error {
	skip := int64(offset)
	limit := int64(length)
	options := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}

	return mg.pageQuery.SimpleFind(data, mg.filter, options)
}
