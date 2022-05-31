package adapter

import (
	"context"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/KejawenLab/bima/v2/paginations"
	"github.com/kamva/mgm/v3"
	paginator "github.com/vcraescu/go-paginator/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MongodbAdapter struct {
		Env        *configs.Env
		Dispatcher *events.Dispatcher
	}

	mongodbPaginator struct {
		context    context.Context
		index      string
		pageQuery  *mgm.Collection
		totalQuery *mgm.Collection
	}
)

func (mg *MongodbAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	model, ok := paginator.Model.(mgm.Model)
	if !ok {
		return nil
	}

	query := mgm.Coll(model)
	event := events.MongodbPagination{
		Query:   query,
		Filters: paginator.Filters,
	}
	mg.Dispatcher.Dispatch(events.PAGINATION_EVENT, &event)

	return newMongodbPaginator(ctx, event.Query)
}

func newMongodbPaginator(context context.Context, query *mgm.Collection) paginator.Adapter {
	var totalQuery *mgm.Collection
	*totalQuery = *query

	return &mongodbPaginator{
		context:    context,
		pageQuery:  query,
		totalQuery: totalQuery,
	}
}

func (mg *mongodbPaginator) Nums() (int64, error) {
	return mg.totalQuery.CountDocuments(mg.context, bson.D{})
}

func (mg *mongodbPaginator) Slice(offset int, length int, data interface{}) error {
	skip := int64(offset)
	limit := int64(length)
	options := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}

	mg.pageQuery.SimpleFind(data, bson.D{}, options)

	return nil
}
