package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	configs "github.com/KejawenLab/bima/v2/configs"
	events "github.com/KejawenLab/bima/v2/events"
	paginations "github.com/KejawenLab/bima/v2/paginations"
)

type Handler struct {
	Env        *configs.Env
	Context    context.Context
	Dispatcher *events.Dispatcher
	Logger     *Logger
	Repository configs.Repository
	Adapter    paginations.Adapter
}

func (h *Handler) Paginate(paginator paginations.Pagination) (paginations.Metadata, []map[string]interface{}) {
	var result []map[string]interface{}
	adapter := h.Adapter.CreateAdapter(h.Context, paginator)
	paginator.Paginate(adapter)
	paginator.Pager.Results(&result)
	next := paginator.Page + 1
	total, _ := paginator.Pager.Nums()

	if paginator.Page*paginator.Limit > int(total) {
		next = -1
	}

	return paginations.Metadata{
		Record:   len(result),
		Page:     paginator.Page,
		Previous: paginator.Page - 1,
		Next:     next,
		Limit:    paginator.Limit,
		Total:    int(total),
	}, result
}

func (h *Handler) Create(v interface{}) error {
	h.Repository.StartTransaction()
	h.Dispatcher.Dispatch(events.BEFORE_CREATE_EVENT, &events.Model{
		Data:       v,
		Repository: h.Repository,
	})

	err := h.Repository.Create(v)
	if err != nil {
		h.Logger.Error("Error when creating resource(s), Rolling back")
		h.Repository.Rollback()

		return err
	}

	if h.Env.Debug {
		m, _ := json.Marshal(v)
		h.Logger.Info(fmt.Sprintf("Elasticsearch result: %s", string(m)))
	}

	h.Dispatcher.Dispatch(events.AFTER_CREATE_EVENT, &events.Model{
		Data:       v,
		Repository: h.Repository,
	})
	h.Repository.Commit()

	return nil
}

func (h *Handler) Update(v interface{}, id string) error {
	h.Repository.StartTransaction()
	h.Dispatcher.Dispatch(events.BEFORE_UPDATE_EVENT, &events.Model{
		Id:         id,
		Data:       v,
		Repository: h.Repository,
	})

	err := h.Repository.Update(v)
	if err != nil {
		h.Logger.Error("Error when creating resource(s), Rolling back")
		h.Repository.Rollback()

		return err
	}

	if h.Env.Debug {
		m, _ := json.Marshal(v)
		h.Logger.Info(fmt.Sprintf("Elasticsearch result: %s", string(m)))
	}

	h.Dispatcher.Dispatch(events.AFTER_UPDATE_EVENT, &events.Model{
		Id:         id,
		Data:       v,
		Repository: h.Repository,
	})
	h.Repository.Commit()

	return nil
}

func (h *Handler) Bind(v interface{}, id string) error {
	return h.Repository.Bind(v, id)
}

func (h *Handler) All(v interface{}) error {
	return h.Repository.All(v)
}

func (h *Handler) Delete(v interface{}, id string) error {
	h.Repository.StartTransaction()
	h.Dispatcher.Dispatch(events.BEFORE_DELETE_EVENT, &events.Model{
		Id:         id,
		Data:       v,
		Repository: h.Repository,
	})

	err := h.Repository.Delete(v, id)
	if err != nil {
		h.Logger.Error("Error when creating resource(s), Rolling back")
		h.Repository.Rollback()

		return err
	}

	if h.Env.Debug {
		m, _ := json.Marshal(v)
		h.Logger.Info(fmt.Sprintf("Delete resources: %s", string(m)))
	}

	h.Dispatcher.Dispatch(events.AFTER_DELETE_EVENT, &events.Model{
		Id:         id,
		Data:       v,
		Repository: h.Repository,
	})
	h.Repository.Commit()

	return nil
}
