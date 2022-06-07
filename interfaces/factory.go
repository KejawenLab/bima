package interfaces

import (
	"sort"
	"time"

	"github.com/KejawenLab/bima/v2/configs"
)

type Application struct {
	Applications []configs.Application
}

func (a *Application) Run(servers []configs.Server) {
	sort.Slice(a.Applications, func(i int, j int) bool {
		return a.Applications[i].Priority() > a.Applications[j].Priority()
	})

	for _, application := range a.Applications {
		time.Sleep(100 * time.Millisecond)
		if application.IsBackground() {
			go application.Run(servers)
		} else {
			application.Run(servers)
		}
	}
}
