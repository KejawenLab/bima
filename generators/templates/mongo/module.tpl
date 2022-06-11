package {{.ModulePluralLowercase}}

import (
	"context"
	"fmt"
	"net/http"

    "github.com/KejawenLab/bima/v2"
	"github.com/jinzhu/copier"
	"{{.PackageName}}/protos/builds"
	"{{.PackageName}}/{{.ModulePluralLowercase}}/models"
	"{{.PackageName}}/{{.ModulePluralLowercase}}/validations"
    "gopkg.in/mgo.v2/bson"
)

type Module struct {
    *bima.Module
	Validator *validations.{{.Module}}
    grpcs.Unimplemented{{.Module}}sServer
}

func (m *Module) GetPaginated(_ context.Context, r *grpcs.Pagination) (*grpcs.{{.Module}}PaginatedResponse, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))
	records := []*grpcs.{{.Module}}{}
	model := models.{{.Module}}{}
	m.Paginator.Model = &model
	m.Paginator.Table = model.CollectionName()

    copier.Copy(m.Request, r)
	m.Paginator.Handle(m.Request)

	metadata, result := m.Handler.Paginate(*m.Paginator)
	for _, v := range result {
	    record := &grpcs.{{.Module}}{}
		data, _ := bson.Marshal(v)
		bson.Unmarshal(data, &model)
		copier.Copy(record, &model)

		record.Id = model.ID.Hex()
		records = append(records, record)
	}

	return &grpcs.{{.Module}}PaginatedResponse{
		Code: http.StatusOK,
		Data: records,
		Meta: &grpcs.PaginationMetadata{
			Record:   int32(metadata.Record),
			Page:     int32(metadata.Page),
			Previous: int32(metadata.Previous),
			Next:     int32(metadata.Next),
			Limit:    int32(metadata.Limit),
			Total:    int32(metadata.Total),
		},
	}, nil
}

func (m *Module) Create(_ context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))

	v := models.{{.Module}}{}
	copier.Copy(&v, r)

	ok, err := m.Validator.Validate(&v)
	if !ok {
		m.Logger.Info(fmt.Sprintf("%+v", err))
		return &grpcs.{{.Module}}Response{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	err = m.Handler.Create(&v)
	if err != nil {
		return &grpcs.{{.Module}}Response{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	r.Id = v.ID.Hex()

	return &grpcs.{{.Module}}Response{
		Code: http.StatusCreated,
		Data: r,
	}, nil
}

func (m *Module) Update(_ context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))

	v := models.{{.Module}}{}
    hold := v
	copier.Copy(&v, r)

	ok, err := m.Validator.Validate(&v)
	if !ok {
		m.Logger.Info(fmt.Sprintf("%+v", err))
		return &grpcs.{{.Module}}Response{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	err = m.Handler.Bind(&hold, r.Id)
	if err != nil {
		m.Logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

		return &grpcs.{{.Module}}Response{
			Code:    http.StatusNotFound,
			Data:    nil,
			Message: err.Error(),
		}, nil
	}

    v.SetID(hold.GetID())
	v.CreatedAt = hold.CreatedAt
	err = m.Handler.Update(&v, v.ID.Hex())
	if err != nil {
		return &grpcs.{{.Module}}Response{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}
    m.Cache.Invalidate(r.Id)

	return &grpcs.{{.Module}}Response{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (m *Module) Get(_ context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))

	var v models.{{.Module}}

	data, found := m.Cache.Get(r.Id)
	if found {
		v = data.(models.{{.Module}})
	} else {
		err := m.Handler.Bind(&v, r.Id)
		if err != nil {
			m.Logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

			return &grpcs.{{.Module}}Response{
				Code:    http.StatusNotFound,
				Data:    nil,
				Message: err.Error(),
			}, nil
		}

		m.Cache.Set(r.Id, v)
	}

	copier.Copy(r, &v)

	return &grpcs.{{.Module}}Response{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (m *Module) Delete(_ context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))

	v := models.{{.Module}}{}

	err := m.Handler.Bind(&v, r.Id)
	if err != nil {
		m.Logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

		return &grpcs.{{.Module}}Response{
			Code:    http.StatusNotFound,
			Data:    nil,
			Message: err.Error(),
		}, nil
	}

    m.Handler.Delete(&v, r.Id)
    m.Cache.Invalidate(r.Id)

	return &grpcs.{{.Module}}Response{
		Code: http.StatusNoContent,
		Data: nil,
	}, nil
}
