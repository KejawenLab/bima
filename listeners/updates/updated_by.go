package updates

import (
	"time"

	configs "github.com/Kejawenlab/bima/v2/configs"
	events "github.com/Kejawenlab/bima/v2/events"
)

type UpdatedBy struct {
	Env *configs.Env
}

func (c *UpdatedBy) Handle(event interface{}) {
	e := event.(*events.Model)
	data := e.Data.(configs.Model)
	data.SetUpdatedBy(c.Env.User)
	data.SetUpdatedAt(time.Now())
	e.Repository.OverrideData(data)
}

func (u *UpdatedBy) Listen() string {
	return events.BEFORE_UPDATE_EVENT
}

func (c *UpdatedBy) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
