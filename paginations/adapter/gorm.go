package adapter

import (
	"context"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/KejawenLab/bima/v2/paginations"
	"github.com/vcraescu/go-paginator/v2"
	"github.com/vcraescu/go-paginator/v2/adapter"
	"gorm.io/gorm"
)

type GormAdapter struct {
	Env        *configs.Env
	Dispatcher *events.Dispatcher
	Database   *gorm.DB
}

func (g *GormAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	query := g.Database.Model(paginator.Model)
	event := events.GormPagination{
		Query:   query,
		Filters: paginator.Filters,
	}
	g.Dispatcher.Dispatch(events.PAGINATION_EVENT, &event)

	return adapter.NewGORMAdapter(event.Query)
}
