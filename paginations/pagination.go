package paginations

import (
	"strings"

	grpcs "github.com/crowdeco/bima/protos/builds"
	paginator "github.com/vcraescu/go-paginator/v2"
)

type (
	Filter struct {
		Field string
		Value string
	}

	Pagination struct {
		Limit      int
		Page       int
		UseCounter bool
		Counter    uint64
		Filters    []Filter
		Search     string
		Pager      paginator.Paginator
		Model      string
	}

	PaginationMeta struct {
		Record   int
		Page     int
		Previous int
		Next     int
		Limit    int
		Total    int
	}
)

func (p *Pagination) Handle(pagination *grpcs.Pagination) {
	if 0 == pagination.Page {
		pagination.Page = 1
	}

	if 0 == pagination.Limit {
		pagination.Limit = 17
	}

	p.Filters = nil
	if len(pagination.Fields) == len(pagination.Values) {
		for k, v := range pagination.Fields {
			if v != "" {
				p.Filters = append(p.Filters, Filter{Field: strings.Title(v), Value: pagination.Values[k]})
			}
		}
	}

	if pagination.Counter > 0 {
		p.UseCounter = true
		p.Counter = pagination.Counter
	}

	p.Limit = int(pagination.Limit)
	p.Page = int(pagination.Page)
}

func (p *Pagination) Paginate(adapter paginator.Adapter) *Pagination {
	pager := paginator.New(adapter, p.Limit)
	pager.SetPage(p.Page)
	p.Pager = pager

	return p
}
