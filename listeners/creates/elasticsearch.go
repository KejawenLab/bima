package creates

import (
	"context"
	"encoding/json"
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

func (c *Elasticsearch) Handle(event interface{}) {
	e := event.(*events.Model)

	m := e.Data.(configs.Model)
	data, _ := json.Marshal(e.Data)

	c.Logger.Info(fmt.Sprintf("Sending data to elasticsearch: %s", string(data)))
	c.Elasticsearch.Index().Index(fmt.Sprintf("%s_%s", c.Env.ServiceCanonicalName, m.TableName())).BodyJson(string(data)).Do(c.Context)
}

func (u *Elasticsearch) Listen() string {
	return events.AFTER_CREATE_EVENT
}

func (c *Elasticsearch) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
