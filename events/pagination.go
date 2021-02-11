package events

import (
	paginations "github.com/crowdeco/bima/paginations"
	services "github.com/crowdeco/bima/services"
	elastic "github.com/olivere/elastic/v7"
)

type Pagination struct {
	Repository *services.Repository
	Query      *elastic.BoolQuery
	Filters    []paginations.Filter
}
