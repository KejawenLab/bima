package paginations

import (
	"fmt"

	configs "github.com/KejawenLab/bima/v2/configs"
	events "github.com/KejawenLab/bima/v2/events"
	"github.com/olivere/elastic/v7"
)

type ElasticsearchFilter struct {
}

func (u *ElasticsearchFilter) Handle(event interface{}) {
	e, ok := event.(*events.ElasticsearchPagination)
	if !ok {
		return
	}

	query := e.Query
	filters := e.Filters
	for _, v := range filters {
		q := elastic.NewWildcardQuery(v.Field, fmt.Sprintf("*%s*", v.Value))
		q.Boost(1.0)
		query.Must(q)
	}
}

func (u *ElasticsearchFilter) Listen() string {
	return events.PAGINATION_EVENT
}

func (u *ElasticsearchFilter) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
