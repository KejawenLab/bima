package updates

import (
	"time"

	configs "github.com/KejawenLab/bima/v2/configs"
	events "github.com/KejawenLab/bima/v2/events"
)

type UpdatedBy struct {
	Env *configs.Env
}

func (c *UpdatedBy) Handle(event interface{}) interface{} {
	e := event.(*events.Model)
	data := e.Data.(configs.Model)

	data.SetUpdatedBy(c.Env.User)
	data.SetUpdatedAt(time.Now())
	e.Repository.OverrideData(data)

	return e
}

func (u *UpdatedBy) Listen() string {
	return events.BEFORE_UPDATE_EVENT
}

func (c *UpdatedBy) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
