package paginations

import (
	"fmt"

	"github.com/KejawenLab/bima/v3/events"
)

type GormFilter struct {
}

func (p *GormFilter) Handle(event interface{}) interface{} {
	e, ok := event.(*events.GormPagination)
	if !ok {
		return event
	}

	query := e.Query
	filters := e.Filters
	for _, v := range filters {
		query.Where(fmt.Sprintf("%s LIKE ?", v.Field), fmt.Sprintf("%%%s%%", v.Value))
	}

	return e
}

func (p *GormFilter) Listen() string {
	return events.PaginationEvent.String()
}

func (p *GormFilter) Priority() int {
	return 255
}
