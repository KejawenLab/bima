package handlers

import (
	"context"

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
	var result []map[string]interface{}
	adapter := h.Adapter.CreateAdapter(h.Context, paginator)
	if adapter == nil {
		h.Logger.Error("Error when creating adapter")

		return paginations.Metadata{}, []map[string]interface{}{}
	}

	var total64 int64
	paginator.Paginate(adapter, &result, &total64)

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
		h.Dispatcher.Dispatch(events.BEFORE_CREATE_EVENT, &events.Model{
			Data:       v,
			Repository: r,
		})

		err := r.Create(v)
		if err != nil {
			h.Logger.Error("Error when creating resource(s), Rolling back")

			return err
		}

		h.Dispatcher.Dispatch(events.AFTER_CREATE_EVENT, &events.Model{
			Data:       v,
			Repository: r,
		})

		return nil
	})
}

func (h *Handler) Update(v interface{}, id string) error {
	return h.Repository.Transaction(func(r configs.Repository) error {
		h.Dispatcher.Dispatch(events.BEFORE_UPDATE_EVENT, &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		err := r.Update(v)
		if err != nil {
			h.Logger.Error("Error when updating resource(s), Rolling back")

			return err
		}

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
		h.Dispatcher.Dispatch(events.BEFORE_DELETE_EVENT, &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		err := r.Delete(v, id)
		if err != nil {
			h.Logger.Error("Error when deleting resource(s), Rolling back")

			return err
		}

		h.Dispatcher.Dispatch(events.AFTER_DELETE_EVENT, &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		return nil
	})
}
