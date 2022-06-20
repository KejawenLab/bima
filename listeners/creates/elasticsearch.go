package creates

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/KejawenLab/bima/v2"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/KejawenLab/bima/v2/models"
	"github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Service       string
	Elasticsearch *elastic.Client
}

func (c *Elasticsearch) Handle(event interface{}) interface{} {
	e := event.(*events.Model)
	m := e.Data.(models.GormModel)

	result := make(chan error)
	go func(r chan<- error) {
		data, _ := json.Marshal(e.Data)

		_, err := c.Elasticsearch.Index().Index(fmt.Sprintf("%s_%s", c.Service, m.TableName())).BodyJson(string(data)).Do(context.Background())

		r <- err
	}(result)

	go func(r <-chan error) {
		if <-r == nil {
			m.SetSyncedAt(time.Now())
			e.Repository.Update(m)
		}
	}(result)

	return e
}

func (u *Elasticsearch) Listen() string {
	return events.AfterCreateEvent.String()
}

func (c *Elasticsearch) Priority() int {
	return bima.HighestPriority + 1
}
