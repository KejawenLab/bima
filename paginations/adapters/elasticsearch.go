package adapters

import (
	"context"
	"strings"

	"github.com/goccy/go-json"

	"github.com/KejawenLab/bima/v4/events"
	"github.com/KejawenLab/bima/v4/loggers"
	"github.com/KejawenLab/bima/v4/paginations"
	"github.com/olivere/elastic/v7"
	"github.com/vcraescu/go-paginator/v2"
)

type (
	ElasticsearchAdapter struct {
		Debug      bool
		Service    string
		Client     *elastic.Client
		Dispatcher *events.Dispatcher
	}

	elasticsearchPaginator struct {
		context    context.Context
		client     *elastic.Client
		index      string
		model      interface{}
		pageQuery  *elastic.BoolQuery
		totalQuery *elastic.BoolQuery
	}
)

func (es *ElasticsearchAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	if es.Client == nil {
		loggers.Logger.Error(ctx, "adapter not configured properly")

		return nil
	}

	query := elastic.NewBoolQuery()
	event := events.ElasticsearchPagination{
		Query:   query,
		Filters: paginator.Filters,
	}

	if es.Debug {
		var log strings.Builder
		log.WriteString("dispatching ")
		log.WriteString(events.PaginationEvent.String())

		loggers.Logger.Debug(ctx, log.String())
	}

	var index strings.Builder

	index.WriteString(es.Service)
	index.WriteString("_")
	index.WriteString(paginator.Table)

	es.Dispatcher.Dispatch(events.PaginationEvent.String(), &event)

	return newElasticsearchPaginator(ctx, es.Client, index.String(), paginator.Model, event.Query)
}

func newElasticsearchPaginator(context context.Context, client *elastic.Client, index string, model interface{}, query *elastic.BoolQuery) paginator.Adapter {
	totalQuery := query
	paginator := elasticsearchPaginator{
		context:    context,
		client:     client,
		index:      index,
		model:      model,
		pageQuery:  query,
		totalQuery: totalQuery,
	}

	return &paginator
}

func (es *elasticsearchPaginator) Nums() (int64, error) {
	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.totalQuery).Do(es.context)
	if err != nil {
		return 0, err
	}

	return result.TotalHits(), nil
}

func (es *elasticsearchPaginator) Slice(offset int, length int, data interface{}) error {
	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.pageQuery).From(offset).Size(length).Do(es.context)
	if err != nil {
		return err
	}

	if result.Hits == nil {
		return nil
	}

	records := make([]map[string]interface{}, 0, result.TotalHits())
	var record map[string]interface{}
	for k, hit := range result.Hits.Hits {
		json.Unmarshal(hit.Source, &record)
		records[k] = record
	}

	temp, _ := json.Marshal(records)

	return json.Unmarshal(temp, data)
}
