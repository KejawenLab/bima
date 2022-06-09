package updates

import (
	"context"
	"encoding/json"
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

func (u *Elasticsearch) Handle(event interface{}) interface{} {
	e := event.(*events.Model)
	m := e.Data.(configs.Model)

	query := elastic.NewMatchQuery("Id", e.Id)

	result, _ := u.Elasticsearch.Search().Index(fmt.Sprintf("%s_%s", u.Service.ConnonicalName, m.TableName())).Query(query).Do(u.Context)
	for _, hit := range result.Hits.Hits {
		u.Elasticsearch.Delete().Index(fmt.Sprintf("%s_%s", u.Service.ConnonicalName, m.TableName())).Id(hit.Id).Do(u.Context)
	}

	data, _ := json.Marshal(e.Data)

	u.Elasticsearch.Index().Index(fmt.Sprintf("%s_%s", u.Service.ConnonicalName, m.TableName())).BodyJson(string(data)).Do(u.Context)

	m.SetSyncedAt(time.Now())
	e.Repository.Update(m)

	return e
}

func (u *Elasticsearch) Listen() string {
	return events.AFTER_UPDATE_EVENT
}

func (u *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
