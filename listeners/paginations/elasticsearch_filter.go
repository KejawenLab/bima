package paginations

import (
	"fmt"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/olivere/elastic/v7"
)

type ElasticsearchFilter struct {
}

func (u *ElasticsearchFilter) Handle(event interface{}) interface{} {
	e, ok := event.(*events.ElasticsearchPagination)
	if !ok {
		return event
	}

	query := e.Query
	filters := e.Filters
	for _, v := range filters {
		q := elastic.NewWildcardQuery(v.Field, fmt.Sprintf("*%s*", v.Value))
		q.Boost(1.0)
		query.Must(q)
	}

	return e
}

func (u *ElasticsearchFilter) Listen() string {
	return events.PaginationEvent.String()
}

func (u *ElasticsearchFilter) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
