package creates

import (
	"time"

	configs "github.com/crowdeco/bima/configs"
	events "github.com/crowdeco/bima/events"
)

type CreatedBy struct {
	Env *configs.Env
}

func (c *CreatedBy) Handle(event interface{}) {
	e := event.(*events.Model)
	data := e.Data.(configs.Model)
	data.SetCreatedBy(c.Env.User)
	data.SetCreatedAt(time.Now())
	e.Repository.OverrideData(data)
}

func (u *CreatedBy) Listen() string {
	return events.BEFORE_CREATE_EVENT
}

func (c *CreatedBy) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
