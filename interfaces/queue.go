package interfaces

import (
	"github.com/KejawenLab/bima/v3/configs"
	"github.com/KejawenLab/bima/v3/messengers"
)

type Queue struct {
	Messenger *messengers.Messenger
}

func (q *Queue) Run(servers []configs.Server) {
	if q.Messenger == nil {
		return
	}

	for _, server := range servers {
		go server.Consume(q.Messenger)
	}
}

func (q *Queue) IsBackground() bool {
	return true
}

func (q *Queue) Priority() int {
	return 0
}
