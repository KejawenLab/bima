package handlers

import (
	"context"
	"fmt"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/KejawenLab/bima/v2/paginations"
)

type Handler struct {
	Context    context.Context
	Dispatcher *events.Dispatcher
	Logger     *Logger
	Repository configs.Repository
	Adapter    paginations.Adapter
}

func (h *Handler) Paginate(paginator paginations.Pagination) (paginations.Metadata, []map[string]interface{}) {
	ctx := context.WithValue(context.Background(), "scope", "handler")

	var result []map[string]interface{}
	adapter := h.Adapter.CreateAdapter(h.Context, paginator)
	if adapter == nil {
		h.Logger.Error(ctx, "Error when creating adapter")

		return paginations.Metadata{}, []map[string]interface{}{}
	}

	var total64 int64
	paginator.Paginate(adapter, &result, &total64)

	h.Logger.Info(ctx, fmt.Sprintf("Total result: %d", total64))

	var total = int(total64)
	next := paginator.Page + 1
	if paginator.Page*paginator.Limit > int(total) {
		next = -1
	}

	return paginations.Metadata{
		Record:   len(result),
		Page:     paginator.Page,
		Previous: paginator.Page - 1,
		Next:     next,
		Limit:    paginator.Limit,
		Total:    total,
	}, result
}

func (h *Handler) Create(v interface{}) error {
	return h.Repository.Transaction(func(r configs.Repository) error {
		ctx := context.WithValue(context.Background(), "scope", "handler")

		h.Logger.Info(ctx, "Dispatching before create event")
		h.Dispatcher.Dispatch(events.BEFORE_CREATE_EVENT, &events.Model{
			Data:       v,
			Repository: r,
		})

		if err := r.Create(v); err != nil {
			h.Logger.Error(ctx, "Error when creating resource(s), Rolling back")

			return err
		}

		h.Logger.Info(ctx, "Dispatching after create event")
		h.Dispatcher.Dispatch(events.AFTER_CREATE_EVENT, &events.Model{
			Data:       v,
			Repository: r,
		})

		return nil
	})
}

func (h *Handler) Update(v interface{}, id string) error {
	return h.Repository.Transaction(func(r configs.Repository) error {
		ctx := context.WithValue(context.Background(), "scope", "handler")

		h.Logger.Info(ctx, "Dispatching before update event")
		h.Dispatcher.Dispatch(events.BEFORE_UPDATE_EVENT, &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		if err := r.Update(v); err != nil {
			h.Logger.Error(ctx, "Error when updating resource(s), Rolling back")

			return err
		}

		h.Logger.Info(ctx, "Dispatching after update event")
		h.Dispatcher.Dispatch(events.AFTER_UPDATE_EVENT, &events.Model{
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

	return h.Repository.Transaction(func(r configs.Repository) error {
		ctx := context.WithValue(context.Background(), "scope", "handler")

		h.Logger.Info(ctx, "Dispatching before delete event")
		h.Dispatcher.Dispatch(events.BEFORE_DELETE_EVENT, &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		if err := r.Delete(v, id); err != nil {
			h.Logger.Error(ctx, "Error when deleting resource(s), Rolling back")

			return err
		}

		h.Logger.Info(ctx, "Dispatching after delete event")
		h.Dispatcher.Dispatch(events.AFTER_DELETE_EVENT, &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		return nil
	})
}
