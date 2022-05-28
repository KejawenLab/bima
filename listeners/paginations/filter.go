package paginations

import (
	"fmt"

	configs "github.com/KejawenLab/bima/v2/configs"
	events "github.com/KejawenLab/bima/v2/events"
	"github.com/olivere/elastic/v7"
)

type (
	ElasticsearchFilter struct {
	}

	GormFilter struct {
	}
)

func (u *GormFilter) Handle(event interface{}) {
	e, ok := event.(*events.GormPagination)
	if !ok {
		return
	}

	query := e.Query
	filters := e.Filters
	for _, v := range filters {
		query.Where(fmt.Sprintf("%s LIKE ?", v.Field), fmt.Sprintf("%%%s%%", v.Value))
	}
}

func (u *GormFilter) Listen() string {
	return events.PAGINATION_EVENT
}

func (u *GormFilter) Priority() int {
	return configs.HIGEST_PRIORITY + 1
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
