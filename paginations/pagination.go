package paginations

import (
	"context"
	"strings"

	paginator "github.com/vcraescu/go-paginator/v2"
)

type (
	Adapter interface {
		CreateAdapter(ctx context.Context, paginator Pagination) paginator.Adapter
	}

	Pagination struct {
		Limit   int
		Page    int
		Filters []Filter
		Search  string
		Pager   paginator.Paginator
		Model   interface{}
		Table   string
	}

	Filter struct {
		Field string
		Value string
	}

	Metadata struct {
		Record   int
		Page     int
		Previous int
		Next     int
		Limit    int
		Total    int
	}

	Request struct {
		Page   int32
		Limit  int32
		Fields []string
		Values []string
	}
)

func (p *Pagination) Handle(request *Request) {
	if 0 == request.Page {
		request.Page = 1
	}

	if 0 == request.Limit {
		request.Limit = 17
	}

	p.Filters = nil
	if len(request.Fields) == len(request.Values) {
		for k, v := range request.Fields {
			if v != "" {
				p.Filters = append(p.Filters, Filter{Field: strings.Title(v), Value: request.Values[k]})
			}
		}
	}

	p.Limit = int(request.Limit)
	p.Page = int(request.Page)
}

func (p *Pagination) Paginate(adapter paginator.Adapter) *Pagination {
	pager := paginator.New(adapter, p.Limit)
	pager.SetPage(p.Page)
	p.Pager = pager

	return p
}
