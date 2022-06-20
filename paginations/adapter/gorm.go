package adapter

import (
	"context"
	"fmt"

	"github.com/KejawenLab/bima/v2/events"
	"github.com/KejawenLab/bima/v2/loggers"
	"github.com/KejawenLab/bima/v2/paginations"
	"github.com/vcraescu/go-paginator/v2"
	"github.com/vcraescu/go-paginator/v2/adapter"
	"gorm.io/gorm"
)

type GormAdapter struct {
	Logger     *loggers.Logger
	Dispatcher *events.Dispatcher
	Database   *gorm.DB
}

func (g *GormAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	query := g.Database.Model(paginator.Model)
	event := events.GormPagination{
		Query:   query,
		Filters: paginator.Filters,
	}

	g.Logger.Debug(ctx, fmt.Sprintf("Dispatching %s", events.PaginationEvent))
	g.Dispatcher.Dispatch(events.PaginationEvent.String(), &event)

	return adapter.NewGORMAdapter(event.Query)
}
