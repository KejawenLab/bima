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
	Model     *{{.Module}}
    grpcs.Unimplemented{{.Module}}sServer
}

func (m *Module) GetPaginated(ctx context.Context, r *grpcs.Pagination) (*grpcs.{{.Module}}PaginatedResponse, error) {
	reqeust := paginations.Request{}

	m.Paginator.Model = *m.Model
	m.Paginator.Table = m.Model.TableName()

    copier.Copy(&reqeust, r)
	m.Paginator.Handle(reqeust)

    records := make([]*grpcs.{{.Module}}, m.Paginator.Limit)
	metadata := m.Handler.Paginate(*m.Paginator, &records)

	return &grpcs.{{.Module}}PaginatedResponse{
		Data: records,
		Meta: &grpcs.PaginationMetadata{
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
	v := m.Model
	copier.Copy(v, r)

	if message, err := utils.Validate(v); err != nil {
		loggers.Logger.Error(ctx, message)

		return nil, status.Error(codes.InvalidArgument, message)
	}

	if err := m.Handler.Create(v); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	r.Id = v.Id

	return r, nil
}

func (m *Module) Update(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	v := m.Model
    hold := *v
	copier.Copy(v, r)

	if message, err := utils.Validate(v); err != nil {
		loggers.Logger.Error(ctx, message)

		return nil, status.Error(codes.InvalidArgument, message)
	}

	if err := m.Handler.Bind(&hold, r.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
	}

    v.Id = r.Id
	v.SetCreatedBy(hold.CreatedBy.String)
	v.SetCreatedAt(hold.CreatedAt.Time)
	if err := m.Handler.Update(v, v.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, "Internal server error")
	}

    m.Cache.Invalidate(r.Id)

	return r, nil
}

func (m *Module) Get(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	v := *m.Model
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
	v := m.Model
	if err := m.Handler.Bind(v, r.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
	}

    m.Handler.Delete(v, r.Id)
    m.Cache.Invalidate(r.Id)

	return nil, nil
}
