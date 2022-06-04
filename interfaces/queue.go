package interfaces

import (
	"fmt"

	configs "github.com/KejawenLab/bima/v2/configs"
	"github.com/fatih/color"
)

type Queue struct {
}

func (q *Queue) Run(servers []configs.Server) {
	if len(servers) > 0 {
		color.New(color.FgCyan, color.Bold).Printf("✓ ")
		fmt.Println("Waiters Ready for Serving Messages...")
	}

	for _, server := range servers {
		go server.RegisterQueueConsumer()
	}
}

func (q *Queue) IsBackground() bool {
	return true
}

func (q *Queue) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
