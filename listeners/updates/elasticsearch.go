package updates

import (
	"context"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/KejawenLab/bima/v3"
	"github.com/KejawenLab/bima/v3/events"
	"github.com/KejawenLab/bima/v3/models"
	"github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Service       string
	Elasticsearch *elastic.Client
}

func (u *Elasticsearch) Handle(event interface{}) interface{} {
	e := event.(*events.Model)
	if u.Elasticsearch == nil {
		return e
	}

	m := e.Data.(models.GormModel)

	var index strings.Builder

	index.WriteString(u.Service)
	index.WriteString("_")
	index.WriteString(m.TableName())

	result := make(chan error)
	go func(r chan<- error) {
		query := elastic.NewMatchQuery("Id", e.Id)

		ctx := context.Background()
		result, _ := u.Elasticsearch.Search().Index(index.String()).Query(query).Do(ctx)
		if result != nil {
			for _, hit := range result.Hits.Hits {
				u.Elasticsearch.Delete().Index(index.String()).Id(hit.Id).Do(ctx)
			}
		}

		data, _ := json.Marshal(e.Data)

		_, err := u.Elasticsearch.Index().Index(index.String()).BodyJson(string(data)).Do(ctx)
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
	return events.AfterUpdateEvent.String()
}

func (u *Elasticsearch) Priority() int {
	return bima.HighestPriority + 1
}
