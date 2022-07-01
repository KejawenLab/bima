package interfaces

import (
	"github.com/KejawenLab/bima/v3/configs"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

func (d *Database) Run(servers []configs.Server) {
	if d.Db == nil {
		return
	}

	for _, server := range servers {
		go server.Migrate(d.Db)
	}
}

func (d *Database) IsBackground() bool {
	return true
}

func (d *Database) Priority() int {
	return 0
}
