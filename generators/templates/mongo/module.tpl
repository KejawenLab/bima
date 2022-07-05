package {{.ModulePluralLowercase}}

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

    "github.com/KejawenLab/bima/v4"
	"github.com/KejawenLab/bima/v4/loggers"
	"github.com/KejawenLab/bima/v4/paginations"
    "github.com/KejawenLab/bima/v4/utils"
	"github.com/goccy/go-json"
	"github.com/jinzhu/copier"
	"{{.PackageName}}/protos/builds"
)

type Module struct {
    *bima.Module
    grpcs.Unimplemented{{.Module}}sServer
}

func (m *Module) GetPaginated(ctx context.Context, r *grpcs.Pagination) (*grpcs.{{.Module}}PaginatedResponse, error) {
	model := {{.Module}}{}
	reqeust := paginations.Request{}

	m.Paginator.Model = &model
	m.Paginator.Table = model.CollectionName()

    copier.Copy(&reqeust, r)
	m.Paginator.Handle(reqeust)

    records := make([]*grpcs.{{.Module}}, 0, m.Paginator.Limit)
	metadata := m.Handler.Paginate(*m.Paginator, &records)

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

func (m *Module) Create(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	v := {{.Module}}{}
	copier.Copy(&v, r)

	if message, err := utils.Validate(&v); err != nil {
		loggers.Logger.Error(ctx, string(message))

		return nil, status.Error(codes.InvalidArgument, string(message))
	}

	if err := m.Handler.Create(&v); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	r.Id = v.ID.Hex()

	return r, nil
}

func (m *Module) Update(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	v := {{.Module}}{}
    hold := v
	copier.Copy(&v, r)

	if message, err := utils.Validate(&v); err != nil {
		loggers.Logger.Error(ctx, string(message))

		return nil, status.Error(codes.InvalidArgument, string(message))
	}

	if err := m.Handler.Bind(&hold, r.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
	}

    v.SetID(hold.GetID())
	v.CreatedAt = hold.CreatedAt
	if err := m.Handler.Update(&v, v.ID.Hex()); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, "Internal server error")
	}

    m.Cache.Invalidate(r.Id)

	return r, nil
}

func (m *Module) Get(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	var v {{.Module}}
	if data, found := m.Cache.Get(r.Id); found {
		err := json.Unmarshal(data, r)
		if err == nil {
			return r, nil
		}
	} else {
		if err := m.Handler.Bind(&v, r.Id); err != nil {
			loggers.Logger.Error(ctx, err.Error())

			return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
		}
	}

	copier.Copy(r, &v)

    data, err := json.Marshal(r)
	if err == nil {
		m.Cache.Set(r.Id, data)
	}

	return r, nil
}

func (m *Module) Delete(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	v := {{.Module}}{}
	if err := m.Handler.Bind(&v, r.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
	}

    m.Handler.Delete(&v, r.Id)
    m.Cache.Invalidate(r.Id)

	return nil, nil
}
