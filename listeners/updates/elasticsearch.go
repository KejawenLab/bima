package updates

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	configs "github.com/KejawenLab/bima/v2/configs"
	events "github.com/KejawenLab/bima/v2/events"
	handlers "github.com/KejawenLab/bima/v2/handlers"
	elastic "github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Env           *configs.Env
	Context       context.Context
	Elasticsearch *elastic.Client
	Logger        *handlers.Logger
}

func (u *Elasticsearch) Handle(event interface{}) {
	e := event.(*events.Model)

	m := e.Data.(configs.Model)
	query := elastic.NewMatchQuery("Id", e.Id)

	u.Logger.Info(fmt.Sprintf("Deleting data in elasticsearch with ID: %s", string(e.Id)))
	result, _ := u.Elasticsearch.Search().Index(fmt.Sprintf("%s_%s", u.Env.Service.ConnonicalName, m.TableName())).Query(query).Do(u.Context)
	for _, hit := range result.Hits.Hits {
		u.Elasticsearch.Delete().Index(fmt.Sprintf("%s_%s", u.Env.Service.ConnonicalName, m.TableName())).Id(hit.Id).Do(u.Context)
	}

	data, _ := json.Marshal(e.Data)

	u.Logger.Info(fmt.Sprintf("Sending data to elasticsearch: %s", string(data)))
	u.Elasticsearch.Index().Index(fmt.Sprintf("%s_%s", u.Env.Service.ConnonicalName, m.TableName())).BodyJson(string(data)).Do(u.Context)

	m.SetSyncedAt(time.Now())
	e.Repository.Update(m)
}

func (u *Elasticsearch) Listen() string {
	return events.AFTER_UPDATE_EVENT
}

func (u *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
