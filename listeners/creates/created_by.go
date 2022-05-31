package creates

import (
	"time"

	configs "github.com/KejawenLab/bima/v2/configs"
	events "github.com/KejawenLab/bima/v2/events"
)

type CreatedBy struct {
	Env *configs.Env
}

func (c *CreatedBy) Handle(event interface{}) interface{} {
	e := event.(*events.Model)
	data := e.Data.(configs.Model)

	data.SetCreatedBy(c.Env.User)
	data.SetCreatedAt(time.Now())
	e.Repository.OverrideData(data)

	return e
}

func (u *CreatedBy) Listen() string {
	return events.BEFORE_CREATE_EVENT
}

func (c *CreatedBy) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
