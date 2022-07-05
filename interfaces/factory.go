package interfaces

import (
	"sort"
	"time"

	"github.com/KejawenLab/bima/v4/configs"
)

type (
	Application interface {
		Run(servers []configs.Server)
		IsBackground() bool
		Priority() int
	}

	Factory struct {
		applications []Application
	}
)

func (f *Factory) Register(applications ...Application) {
	for _, application := range applications {
		f.Add(application)
	}
}

func (f *Factory) Add(application Application) {
	f.applications = append(f.applications, application)
}

func (f *Factory) Run(servers []configs.Server) {
	sort.Slice(f.applications, func(i int, j int) bool {
		return f.applications[i].Priority() > f.applications[j].Priority()
	})

	for _, application := range f.applications {
		if application.IsBackground() {
			go application.Run(servers)
		} else {
			time.Sleep(100 * time.Millisecond)
			application.Run(servers)
		}
	}
}
