package adapters

import (
	"context"
	"strings"

	"github.com/KejawenLab/bima/v4/configs"
	"github.com/KejawenLab/bima/v4/events"
	"github.com/KejawenLab/bima/v4/loggers"
	"github.com/KejawenLab/bima/v4/paginations"
	"github.com/vcraescu/go-paginator/v2"
	"gorm.io/gorm"
)

type (
	GormAdapter struct {
		Debug      bool
		Dispatcher *events.Dispatcher
	}

	gormPaginator struct {
		query *gorm.DB
		total int64
	}
)

func (g *GormAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	if configs.Database == nil {
		loggers.Logger.Error(ctx, "adapter not configured properly")

		return nil
	}

	query := configs.Database.Model(paginator.Model)
	event := events.GormPagination{
		Query:   query,
		Filters: paginator.Filters,
	}

	if g.Debug {
		var log strings.Builder
		log.WriteString("dispatching ")
		log.WriteString(events.PaginationEvent.String())

		loggers.Logger.Debug(ctx, log.String())
	}

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
