package interfaces

import (
	"github.com/KejawenLab/bima/v4/configs"
)

type Database struct {
}

func (d *Database) Run(servers []configs.Server) {
	if configs.Database == nil {
		return
	}

	for _, server := range servers {
		go server.Migrate(configs.Database)
	}
}

func (d *Database) IsBackground() bool {
	return true
}

func (d *Database) Priority() int {
	return 0
}
