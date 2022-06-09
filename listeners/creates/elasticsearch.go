package creates

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

func (c *Elasticsearch) Handle(event interface{}) interface{} {
	e := event.(*events.Model)
	m := e.Data.(configs.Model)
	data, _ := json.Marshal(e.Data)

	c.Elasticsearch.Index().Index(fmt.Sprintf("%s_%s", c.Service.ConnonicalName, m.TableName())).BodyJson(string(data)).Do(c.Context)

	m.SetSyncedAt(time.Now())
	e.Repository.Update(m)

	return e
}

func (u *Elasticsearch) Listen() string {
	return events.AFTER_CREATE_EVENT
}

func (c *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
