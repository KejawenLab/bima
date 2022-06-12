package paginations

import (
	"context"
	"strings"

	"github.com/vcraescu/go-paginator/v2"
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
	if request.Page == 0 {
		request.Page = 1
	}

	if request.Limit == 0 {
		request.Limit = 17
	}

	p.Limit = int(request.Limit)
	p.Page = int(request.Page)

	n := len(request.Fields)
	if n == 0 {
		return
	}

	p.Filters = []Filter{}
	if n != len(request.Values) {
		return
	}

	for k, v := range request.Fields {
		if v != "" {
			p.Filters = append(p.Filters, Filter{Field: strings.Title(v), Value: request.Values[k]})
		}
	}
}

func (p *Pagination) Paginate(adapter paginator.Adapter) *Pagination {
	pager := paginator.New(adapter, p.Limit)
	pager.SetPage(p.Page)
	p.Pager = pager

	return p
}
