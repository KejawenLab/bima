package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/KejawenLab/bima/v2/paginations"
	elastic "github.com/olivere/elastic/v7"
	paginator "github.com/vcraescu/go-paginator/v2"
)

type (
	ElasticsearchAdapter struct {
		Env        *configs.Env
		Client     *elastic.Client
		Dispatcher *events.Dispatcher
	}

	ElasticsearchPaginator struct {
		context    context.Context
		client     *elastic.Client
		index      string
		pageQuery  *elastic.BoolQuery
		totalQuery *elastic.BoolQuery
	}
)

func (es *ElasticsearchAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	query := elastic.NewBoolQuery()
	es.Dispatcher.Dispatch(events.PAGINATION_EVENT, &events.ElasticsearchPagination{
		Query:   query,
		Filters: paginator.Filters,
	})

	return newElasticsearchPaginator(ctx, es.Client, fmt.Sprintf("%s_%s", es.Env.Service.ConnonicalName, paginator.Table), query)
}

func newElasticsearchPaginator(context context.Context, client *elastic.Client, index string, query *elastic.BoolQuery) paginator.Adapter {
	var totalQuery *elastic.BoolQuery
	*totalQuery = *query

	return &ElasticsearchPaginator{
		context:    context,
		client:     client,
		index:      index,
		pageQuery:  query,
		totalQuery: totalQuery,
	}
}

func (es *ElasticsearchPaginator) Nums() (int64, error) {
	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.totalQuery).Do(es.context)
	if err != nil {
		log.Printf("%s", err.Error())
		return 0, nil
	}

	return result.TotalHits(), nil
}

func (es *ElasticsearchPaginator) Slice(offset int, length int, data interface{}) error {
	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.pageQuery).From(offset).Size(length).Do(es.context)
	if err != nil {
		log.Printf("%s", err.Error())
		return nil
	}

	records := data.(*[]interface{})
	var record interface{}
	for _, hit := range result.Hits.Hits {
		json.Unmarshal(hit.Source, &record)

		*records = append(*records, record)
	}

	data = *records

	return nil
}
