package handlers

import (
	"context"
	"fmt"

	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/repositories"
)

type Handler struct {
	Logger     *loggers.Logger
	Dispatcher *events.Dispatcher
	Repository repositories.Repository
	Adapter    paginations.Adapter
}

func (h *Handler) Paginate(paginator paginations.Pagination, result interface{}) paginations.Metadata {
	ctx := context.WithValue(context.Background(), "scope", "handler")

	adapter := h.Adapter.CreateAdapter(ctx, paginator)
	if adapter == nil {
		h.Logger.Error(ctx, "Error when creating adapter")

		return paginations.Metadata{}
	}

	var total64 int64
	paginator.Paginate(adapter, result, &total64)

	h.Logger.Debug(ctx, fmt.Sprintf("Total result: %d", total64))

	var total = int(total64)
	next := paginator.Page + 1
	if paginator.Page*paginator.Limit > int(total) {
		next = -1
	}

	return paginations.Metadata{
		Page:     paginator.Page,
		Previous: paginator.Page - 1,
		Next:     next,
		Limit:    paginator.Limit,
		Total:    total,
	}
}

func (h *Handler) Create(v interface{}) error {
	return h.Repository.Transaction(func(r repositories.Repository) error {
		ctx := context.WithValue(context.Background(), "scope", "handler")

		h.Logger.Debug(ctx, fmt.Sprintf("Dispatching %s", events.BeforeCreateEvent))
		h.Dispatcher.Dispatch(events.BeforeCreateEvent.String(), &events.Model{
			Data:       v,
			Repository: r,
		})

		if err := r.Create(v); err != nil {
			h.Logger.Error(ctx, "Error when creating resource(s), Rolling back")

			return err
		}

		h.Logger.Debug(ctx, fmt.Sprintf("Dispatching %s", events.AfterCreateEvent))
		h.Dispatcher.Dispatch(events.AfterCreateEvent.String(), &events.Model{
			Data:       v,
			Repository: r,
		})

		return nil
	})
}

func (h *Handler) Update(v interface{}, id string) error {
	return h.Repository.Transaction(func(r repositories.Repository) error {
		ctx := context.WithValue(context.Background(), "scope", "handler")

		h.Logger.Debug(ctx, fmt.Sprintf("Dispatching %s", events.BeforeUpdateEvent))
		h.Dispatcher.Dispatch(events.BeforeUpdateEvent.String(), &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		if err := r.Update(v); err != nil {
			h.Logger.Error(ctx, "Error when updating resource(s), Rolling back")

			return err
		}

		h.Logger.Debug(ctx, fmt.Sprintf("Dispatching %s", events.AfterUpdateEvent))
		h.Dispatcher.Dispatch(events.AfterUpdateEvent.String(), &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		return nil
	})
}

func (h *Handler) Bind(v interface{}, id string) error {
	return h.Repository.Bind(v, id)
}

func (h *Handler) All(v interface{}) error {
	return h.Repository.All(v)
}

func (h *Handler) Delete(v interface{}, id string) error {

	return h.Repository.Transaction(func(r repositories.Repository) error {
		ctx := context.WithValue(context.Background(), "scope", "handler")

		h.Logger.Debug(ctx, fmt.Sprintf("Dispatching %s", events.BeforeDeleteEvent))
		h.Dispatcher.Dispatch(events.BeforeDeleteEvent.String(), &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		if err := r.Delete(v, id); err != nil {
			h.Logger.Error(ctx, "Error when deleting resource(s), Rolling back")

			return err
		}

		h.Logger.Debug(ctx, fmt.Sprintf("Dispatching %s", events.AfterDeleteEvent))
		h.Dispatcher.Dispatch(events.AfterDeleteEvent.String(), &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		return nil
	})
}
