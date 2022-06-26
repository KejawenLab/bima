package {{.ModulePluralLowercase}}

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

    "github.com/KejawenLab/bima/v3"
	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/utils"
	"github.com/jinzhu/copier"
	"{{.PackageName}}/protos/builds"
)

type Module struct {
    *bima.Module
	Model     *{{.Module}}
    grpcs.Unimplemented{{.Module}}sServer
}

func (m *Module) GetPaginated(ctx context.Context, r *grpcs.Pagination) (*grpcs.{{.Module}}PaginatedResponse, error) {
	m.Logger.Debug(context.WithValue(ctx, "scope", "{{.ModuleLowercase}}"), fmt.Sprintf("%+v", r))
	records := []*grpcs.{{.Module}}{}
	model := {{.Module}}{}
	reqeust := paginations.Request{}

	m.Paginator.Model = model
	m.Paginator.Table = model.TableName()

    copier.Copy(&reqeust, r)
	m.Paginator.Handle(reqeust)

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

func (m *Module) Create(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	v := m.Model
	copier.Copy(v, r)

	if message, err := utils.Validate(v); err != nil {
		m.Logger.Error(ctx, message)

		return nil, status.Error(codes.InvalidArgument, message)
	}

	if err := m.Handler.Create(v); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	r.Id = v.Id

	return &grpcs.{{.Module}}Response{
		{{.Module}}: r,
	}, nil
}

func (m *Module) Update(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	v := m.Model
    hold := *v
	copier.Copy(v, r)

	if message, err := utils.Validate(v); err != nil {
		m.Logger.Error(ctx, message)

		return nil, status.Error(codes.InvalidArgument, message)
	}

	if err := m.Handler.Bind(&hold, r.Id); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
	}

    v.Id = r.Id
	v.SetCreatedBy(&configs.User{Id: hold.CreatedBy.String})
	v.SetCreatedAt(hold.CreatedAt.Time)
	if err := m.Handler.Update(v, v.Id); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, "Internal server error")
	}

    m.Cache.Invalidate(r.Id)

	return &grpcs.{{.Module}}Response{
		{{.Module}}: r,
	}, nil
}

func (m *Module) Get(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	var v {{.Module}}
	if data, found := m.Cache.Get(r.Id); found {
		v = data.({{.Module}})
	} else {
		if err := m.Handler.Bind(&v, r.Id); err != nil {
			m.Logger.Error(ctx, err.Error())

			return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
		}

		m.Cache.Set(r.Id, v)
	}

	copier.Copy(r, &v)

	return &grpcs.{{.Module}}Response{
		{{.Module}}: r,
	}, nil
}

func (m *Module) Delete(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	v := m.Model
	if err := m.Handler.Bind(v, r.Id); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
	}

    m.Handler.Delete(v, r.Id)
    m.Cache.Invalidate(r.Id)

	return &grpcs.{{.Module}}Response{
		{{.Module}}: nil,
	}, nil
}
