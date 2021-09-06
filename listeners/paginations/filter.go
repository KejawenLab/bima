package paginations

import (
	"fmt"

	configs "github.com/Kejawenlab/bima/v2/configs"
	events "github.com/Kejawenlab/bima/v2/events"
	"github.com/olivere/elastic/v7"
)

type Filter struct {
}

func (u *Filter) Handle(event interface{}) {
	e := event.(*events.Pagination)
	query := e.Query
	filters := e.Filters

	for _, v := range filters {
		q := elastic.NewWildcardQuery(v.Field, fmt.Sprintf("*%s*", v.Value))
		q.Boost(1.0)
		query.Must(q)
	}
}

func (u *Filter) Listen() string {
	return events.PAGINATION_EVENT
}

func (u *Filter) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
