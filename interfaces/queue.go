package interfaces

import "github.com/KejawenLab/bima/v3/configs"

type Queue struct {
}

func (q *Queue) Run(servers []configs.Server) {
	for _, server := range servers {
		go server.RegisterQueueConsumer()
	}
}

func (q *Queue) IsBackground() bool {
	return true
}

func (q *Queue) Priority() int {
	return 0
}
