package interfaces

import (
	"sort"
	"time"

	configs "github.com/crowdeco/bima/configs"
)

type (
	Application struct {
		Applications []configs.Application
	}
)

func (a *Application) Run(servers []configs.Server) {
	sort.Slice(a.Applications, func(i, j int) bool {
		return a.Applications[i].Priority() > a.Applications[j].Priority()
	})

	for _, application := range a.Applications {
		time.Sleep(300 * time.Millisecond)
		if application.IsBackground() {
			go application.Run(servers)
		} else {
			application.Run(servers)
		}
	}
}
