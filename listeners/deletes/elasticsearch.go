package deletes

import (
	"context"
	"fmt"
	"time"

	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Service       configs.Service
	Context       context.Context
	Elasticsearch *elastic.Client
}

func (d *Elasticsearch) Handle(event interface{}) interface{} {
	e := event.(*events.Model)
	m := e.Data.(configs.Model)

	query := elastic.NewMatchQuery("Id", e.Id)

	result, _ := d.Elasticsearch.Search().Index(fmt.Sprintf("%s_%s", d.Service.ConnonicalName, m.TableName())).Query(query).Do(d.Context)
	for _, hit := range result.Hits.Hits {
		d.Elasticsearch.Delete().Index(fmt.Sprintf("%s_%s", d.Service.ConnonicalName, m.TableName())).Id(hit.Id).Do(d.Context)
	}

	m.SetSyncedAt(time.Now())
	e.Repository.Update(m)

	return e
}

func (d *Elasticsearch) Listen() string {
	return events.AFTER_DELETE_EVENT
}

func (d *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
