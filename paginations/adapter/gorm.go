package adapter

import (
	"bytes"
	"context"

	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/vcraescu/go-paginator/v2"
	"gorm.io/gorm"
)

type (
	GormAdapter struct {
		Logger     *loggers.Logger
		Dispatcher *events.Dispatcher
		Database   *gorm.DB
	}

	gormPaginator struct {
		query *gorm.DB
		total int64
	}
)

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

	var total int64
	event.Query.Count(&total)

	return newGormPaginator(event.Query, total)
}

func newGormPaginator(query *gorm.DB, total int64) paginator.Adapter {
	return &gormPaginator{
		query: query,
		total: total,
	}
}

func (gm *gormPaginator) Nums() (int64, error) {
	return gm.total, nil
}

func (gm *gormPaginator) Slice(offset int, length int, data interface{}) error {
	return gm.query.Limit(length).Offset(offset).Find(data).Error
}
