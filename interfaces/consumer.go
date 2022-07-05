package interfaces

import (
	"github.com/KejawenLab/bima/v4/configs"
	"github.com/KejawenLab/bima/v4/messengers"
)

type Consumer struct {
	Messenger *messengers.Messenger
}

func (q *Consumer) Run(servers []configs.Server) {
	if q.Messenger == nil {
		return
	}

	for _, server := range servers {
		go server.Consume(q.Messenger)
	}
}

func (q *Consumer) IsBackground() bool {
	return true
}

func (q *Consumer) Priority() int {
	return 0
}
