package paginations

import (
	"bytes"

	"github.com/KejawenLab/bima/v3/events"
)

type GormFilter struct {
}

func (p *GormFilter) Handle(event interface{}) interface{} {
	e, ok := event.(*events.GormPagination)
	if !ok {
		return event
	}

	var likeClausal bytes.Buffer
	var likeValue bytes.Buffer
	query := e.Query
	filters := e.Filters
	for _, v := range filters {
		likeClausal.Reset()
		likeClausal.WriteString(v.Field)
		likeClausal.WriteString(" LIKE ?")

		likeValue.Reset()
		likeValue.WriteString("%")
		likeValue.WriteString(v.Value)
		likeValue.WriteString("%")

		query.Where(likeClausal.String(), likeValue.String())
	}

	return e
}

func (p *GormFilter) Listen() string {
	return events.PaginationEvent.String()
}

func (p *GormFilter) Priority() int {
	return 255
}
