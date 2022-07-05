package interfaces

import (
	"github.com/KejawenLab/bima/v4/configs"
	"github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	Client *elastic.Client
}

func (e *Elasticsearch) Run(servers []configs.Server) {
	if e.Client == nil {
		return
	}

	for _, server := range servers {
		go server.Sync(e.Client)
	}
}

func (e *Elasticsearch) IsBackground() bool {
	return true
}

func (e *Elasticsearch) Priority() int {
	return 0
}
