package interfaces

import (
	"github.com/KejawenLab/bima/v2/configs"
)

type Elasticsearch struct {
}

func (e *Elasticsearch) Run(servers []configs.Server) {
	for _, server := range servers {
		go server.RepopulateData()
	}
}

func (e *Elasticsearch) IsBackground() bool {
	return true
}

func (e *Elasticsearch) Priority() int {
	return 0
}
