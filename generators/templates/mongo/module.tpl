package {{.ModulePluralLowercase}}

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (m *Module) GetPaginated(ctx context.Context, r *grpcs.Pagination) (*grpcs.{{.Module}}PaginatedResponse, error) {
	m.Logger.Debug(context.WithValue(ctx, "scope", "{{.ModuleLowercase}}"), fmt.Sprintf("%+v", r))
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

func (m *Module) Create(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	v := models.{{.Module}}{}
	copier.Copy(&v, r)

	if ok, err := m.Validator.Validate(&v); !ok {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := m.Handler.Create(&v); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, err.Error())
	}

	r.Id = v.ID.Hex()

	return &grpcs.{{.Module}}Response{
		{{.Module}}: r,
	}, nil
}

func (m *Module) Update(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	v := models.{{.Module}}{}
    hold := v
	copier.Copy(&v, r)

	if ok, err := m.Validator.Validate(&v); !ok {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := m.Handler.Bind(&hold, r.Id); err != nil {
		msg := fmt.Sprintf("Data with ID '%s' not found.", r.Id)
		m.Logger.Error(ctx, msg)

		return nil, status.Error(codes.NotFound, msg)
	}

    v.SetID(hold.GetID())
	v.CreatedAt = hold.CreatedAt
	if err := m.Handler.Update(&v, v.ID.Hex()); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, err.Error())
	}

    m.Cache.Invalidate(r.Id)

	return &grpcs.{{.Module}}Response{
		{{.Module}}: r,
	}, nil
}

func (m *Module) Get(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

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
		{{.Module}}: r,
	}, nil
}

func (m *Module) Delete(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}Response, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	v := models.{{.Module}}{}
	if err := m.Handler.Bind(&v, r.Id); err != nil {
		msg := fmt.Sprintf("Data with ID '%s' not found.", r.Id)
		m.Logger.Error(ctx, msg)

		return nil, status.Error(codes.NotFound, msg)
	}

    m.Handler.Delete(&v, r.Id)
    m.Cache.Invalidate(r.Id)

	return &grpcs.{{.Module}}Response{
		{{.Module}}: nil,
	}, nil
}
