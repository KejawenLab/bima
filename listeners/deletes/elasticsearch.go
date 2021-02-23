package deletes

import (
	"context"
	"fmt"

	configs "github.com/crowdeco/bima/configs"
	events "github.com/crowdeco/bima/events"
	handlers "github.com/crowdeco/bima/handlers"
	elastic "github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Env           *configs.Env
	Context       context.Context
	Elasticsearch *elastic.Client
	Logger        *handlers.Logger
}

func (d *Elasticsearch) Handle(event interface{}) {
	e := event.(*events.Model)

	m := e.Data.(configs.Model)
	query := elastic.NewMatchQuery("Id", e.Id)

	d.Logger.Info(fmt.Sprintf("Deleting data in elasticsearch with ID: %s", string(e.Id)))
	result, _ := d.Elasticsearch.Search().Index(fmt.Sprintf("%s_%s", d.Env.ServiceCanonicalName, m.TableName())).Query(query).Do(d.Context)
	for _, hit := range result.Hits.Hits {
		d.Elasticsearch.Delete().Index(fmt.Sprintf("%s_%s", d.Env.ServiceCanonicalName, m.TableName())).Id(hit.Id).Do(d.Context)
	}
}

func (d *Elasticsearch) Listen() string {
	return events.AFTER_DELETE_EVENT
}

func (d *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
