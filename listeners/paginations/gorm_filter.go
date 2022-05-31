package paginations

import (
	"fmt"

	configs "github.com/KejawenLab/bima/v2/configs"
	events "github.com/KejawenLab/bima/v2/events"
)

type GormFilter struct {
}

func (u *GormFilter) Handle(event interface{}) interface{} {
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

func (u *GormFilter) Listen() string {
	return events.PAGINATION_EVENT
}

func (u *GormFilter) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
