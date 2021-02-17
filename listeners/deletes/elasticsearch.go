package deletes

import (
	"context"
	"fmt"

	configs "github.com/crowdeco/bima/configs"
	events "github.com/crowdeco/bima/events"
	elastic "github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Env           *configs.Env
	Context       context.Context
	Elasticsearch *elastic.Client
}

func (d *Elasticsearch) Handle(event interface{}) {
	e := event.(*events.Model)

	m := e.Data.(configs.Model)
	query := elastic.NewMatchQuery("Id", e.Id)
	result, _ := d.Elasticsearch.Search().Index(fmt.Sprintf("%s_%s", d.Env.ServiceConicalName, m.TableName())).Query(query).Do(d.Context)
	for _, hit := range result.Hits.Hits {
		d.Elasticsearch.Delete().Index(fmt.Sprintf("%s_%s", d.Env.ServiceConicalName, m.TableName())).Id(hit.Id).Do(d.Context)
	}
}

func (d *Elasticsearch) Listen() string {
	return events.AFTER_DELETE_EVENT
}

func (d *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
