package {{.ModulePluralLowercase}}

import (
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

    "github.com/KejawenLab/bima/v2"
	"github.com/KejawenLab/bima/v2/configs"
	"github.com/jinzhu/copier"
	"{{.PackageName}}/protos/builds"
	"{{.PackageName}}/{{.ModulePluralLowercase}}/models"
	"{{.PackageName}}/{{.ModulePluralLowercase}}/validations"
)

type Module struct {
    *bima.Module
	Model     *models.{{.Module}}
	Validator *validations.{{.Module}}
    grpcs.Unimplemented{{.Module}}sServer
}

func (m *Module) GetPaginated(_ context.Context, r *grpcs.Pagination) (*grpcs.{{.Module}}PaginatedResponse, error) {
	m.Logger.Info(context.WithValue(context.Background(), "scope", "{{.ModuleLowercase}}"), fmt.Sprintf("%+v", r))
	records := []*grpcs.{{.Module}}{}
	model := models.{{.Module}}{}
	m.Paginator.Model = model
	m.Paginator.Table = model.TableName()

    copier.Copy(m.Request, r)
	m.Paginator.Handle(m.Request)

	metadata, result := m.Handler.Paginate(*m.Paginator)
	for _, v := range result {
	    record := &grpcs.{{.Module}}{}
		data, _ := json.Marshal(v)
		json.Unmarshal(data, &model)
		copier.Copy(record, &model)

		record.Id = model.Id
		records = append(records, record)
	}

	return &grpcs.{{.Module}}PaginatedResponse{
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
    ctx := context.WithValue(context.Background(), "scope", "{{.ModuleLowercase}}")
	m.Logger.Info(ctx, fmt.Sprintf("%+v", r))

	v := m.Model
	copier.Copy(v, r)

	if ok, err := m.Validator.Validate(v); !ok {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := m.Handler.Create(v); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, err.Error())
	}

	r.Id = v.Id

	return &grpcs.{{.Module}}Response{
		Data: r,
	}, nil
}

func (m *Module) Update(_ context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx := context.WithValue(context.Background(), "scope", "{{.ModuleLowercase}}")
	m.Logger.Info(ctx, fmt.Sprintf("%+v", r))

	v := m.Model
    hold := *v
	copier.Copy(v, r)

	if ok, err := m.Validator.Validate(v); !ok {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := m.Handler.Bind(&hold, r.Id); err != nil {
		msg := fmt.Sprintf("Data with ID '%s' not found.", r.Id)
		m.Logger.Error(ctx, msg)

		return nil, status.Error(codes.NotFound, msg)
	}

    v.Id = r.Id
	v.SetCreatedBy(&configs.User{Id: hold.CreatedBy.String})
	v.SetCreatedAt(hold.CreatedAt.Time)
	if err := m.Handler.Update(v, v.Id); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, err.Error())
	}

    m.Cache.Invalidate(r.Id)

	return &grpcs.{{.Module}}Response{
		Data: r,
	}, nil
}

func (m *Module) Get(_ context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx := context.WithValue(context.Background(), "scope", "{{.ModuleLowercase}}")
	m.Logger.Info(ctx, fmt.Sprintf("%+v", r))

	var v models.{{.Module}}
	if data, found := m.Cache.Get(r.Id); found {
		v = data.(models.{{.Module}})
	} else {
		if err := m.Handler.Bind(&v, r.Id); err != nil {
			msg := fmt.Sprintf("Data with ID '%s' not found.", r.Id)
			m.Logger.Error(ctx, msg)

			return nil, status.Error(codes.NotFound, msg)
		}

		m.Cache.Set(r.Id, v)
	}

	copier.Copy(r, &v)

	return &grpcs.{{.Module}}Response{
		Data: r,
	}, nil
}

func (m *Module) Delete(_ context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx := context.WithValue(context.Background(), "scope", "{{.ModuleLowercase}}")
	m.Logger.Info(ctx, fmt.Sprintf("%+v", r))

	v := m.Model
	if err := m.Handler.Bind(v, r.Id); err != nil {
		msg := fmt.Sprintf("Data with ID '%s' not found.", r.Id)
		m.Logger.Error(ctx, msg)

		return nil, status.Error(codes.NotFound, msg)
	}

    m.Handler.Delete(v, r.Id)
    m.Cache.Invalidate(r.Id)

	return &grpcs.{{.Module}}Response{
		Data: nil,
	}, nil
}
