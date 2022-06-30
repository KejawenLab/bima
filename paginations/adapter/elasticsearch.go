package adapter

import (
	"bytes"
	"context"

	"github.com/goccy/go-json"

	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/paginations"
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
		pageQuery  *elastic.BoolQuery
		totalQuery *elastic.BoolQuery
	}
)

func (es *ElasticsearchAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	if es.Client == nil {
		loggers.Logger.Fatal(ctx, "Adapter not configured properly")

		return nil
	}

	query := elastic.NewBoolQuery()
	event := events.ElasticsearchPagination{
		Query:   query,
		Filters: paginator.Filters,
	}

	if es.Debug {
		var log bytes.Buffer
		log.WriteString("Dispatching ")
		log.WriteString(events.PaginationEvent.String())

		loggers.Logger.Debug(ctx, log.String())
	}

	var index bytes.Buffer

	index.WriteString(es.Service)
	index.WriteString("_")
	index.WriteString(paginator.Table)

	return newElasticsearchPaginator(ctx, es.Client, index.String(), event.Query)
}

func newElasticsearchPaginator(context context.Context, client *elastic.Client, index string, query *elastic.BoolQuery) paginator.Adapter {
	var totalQuery *elastic.BoolQuery
	*totalQuery = *query

	return &elasticsearchPaginator{
		context:    context,
		client:     client,
		index:      index,
		pageQuery:  query,
		totalQuery: totalQuery,
	}
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

	records := data.(*[]map[string]interface{})
	var record map[string]interface{}
	for _, hit := range result.Hits.Hits {
		json.Unmarshal(hit.Source, &record)
		*records = append(*records, record)
	}

	return nil
}
