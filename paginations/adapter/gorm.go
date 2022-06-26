package adapter

import (
	"bytes"
	"context"

	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/paginations"
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

	var log bytes.Buffer
	log.WriteString("Dispatching ")
	log.WriteString(events.PaginationEvent.String())

	g.Logger.Debug(ctx, log.String())
	g.Dispatcher.Dispatch(events.PaginationEvent.String(), &event)

	return adapter.NewGORMAdapter(event.Query)
}
