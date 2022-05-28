package adapter

import (
	"context"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/KejawenLab/bima/v2/paginations"
	"github.com/KejawenLab/bima/v2/services"
	paginator "github.com/vcraescu/go-paginator/v2"
	"github.com/vcraescu/go-paginator/v2/adapter"
)

type GormAdapter struct {
	Env        *configs.Env
	Logger     *handlers.Logger
	Dispatcher *events.Dispatcher
	Repository *services.Repository
}

func (g *GormAdapter) CreateAdapter(ctx context.Context, paginator paginations.Pagination) paginator.Adapter {
	query := g.Repository.Database.Model(paginator.Model)

	g.Dispatcher.Dispatch(events.PAGINATION_EVENT, &events.GormPagination{
		Repository: g.Repository,
		Query:      query,
		Filters:    paginator.Filters,
	})

	return adapter.NewGORMAdapter(query)
}
