package deletes

import (
	"bytes"
	"context"
	"time"

	"github.com/KejawenLab/bima/v3"
	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/models"
	"github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Service       string
	Elasticsearch *elastic.Client
}

func (d *Elasticsearch) Handle(event interface{}) interface{} {
	e := event.(*events.Model)
	if d.Elasticsearch == nil {
		return e
	}

	m := e.Data.(models.GormModel)

	var index bytes.Buffer
	index.WriteString(d.Service)
	index.WriteString("_")
	index.WriteString(m.TableName())

	result := make(chan error)
	go func(c chan<- error) {
		query := elastic.NewMatchQuery("Id", e.Id)

		ctx := context.Background()
		result, _ := d.Elasticsearch.Search().Index(index.String()).Query(query).Do(ctx)
		if result != nil {
			for _, hit := range result.Hits.Hits {
				d.Elasticsearch.Delete().Index(index.String()).Id(hit.Id).Do(ctx)
			}
		}

		c <- nil
	}(result)

	go func(r <-chan error) {
		if <-r == nil {
			m.SetSyncedAt(time.Now())
			e.Repository.Update(m)
		}
	}(result)

	return e
}

func (d *Elasticsearch) Listen() string {
	return events.AfterDeleteEvent.String()
}

func (d *Elasticsearch) Priority() int {
	return bima.HighestPriority + 1
}
