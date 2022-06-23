package interfaces

import (
	"sort"
	"time"

	"github.com/KejawenLab/bima/v3/configs"
)

type (
	Application interface {
		Run(servers []configs.Server)
		IsBackground() bool
		Priority() int
	}

	Factory struct {
		Applications []Application
	}
)

func (a *Factory) Run(servers []configs.Server) {
	sort.Slice(a.Applications, func(i int, j int) bool {
		return a.Applications[i].Priority() > a.Applications[j].Priority()
	})

	for _, application := range a.Applications {
		if application.IsBackground() {
			go application.Run(servers)
		} else {
			time.Sleep(100 * time.Millisecond)
			application.Run(servers)
		}
	}
}
