package paginations

import (
	"strings"

	"github.com/KejawenLab/bima/v4/events"
	"github.com/olivere/elastic/v7"
)

type ElasticsearchFilter struct {
}

func (p *ElasticsearchFilter) Handle(event interface{}) interface{} {
	e, ok := event.(*events.ElasticsearchPagination)
	if !ok {
		return event
	}

	var wildCard strings.Builder
	for _, v := range e.Filters {
		wildCard.Reset()
		wildCard.WriteString("*")
		wildCard.WriteString(v.Value)
		wildCard.WriteString("*")

		q := elastic.NewWildcardQuery(v.Field, wildCard.String())
		q.Boost(1.0)
		e.Query.Must(q)
	}

	return e
}

func (p *ElasticsearchFilter) Listen() string {
	return events.PaginationEvent.String()
}

func (p *ElasticsearchFilter) Priority() int {
	return 255
}
