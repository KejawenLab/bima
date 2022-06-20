package interfaces

import "github.com/KejawenLab/bima/v2/configs"

type Database struct {
}

func (d *Database) Run(servers []configs.Server) {
	for _, server := range servers {
		go server.RegisterAutoMigrate()
	}
}

func (d *Database) IsBackground() bool {
	return true
}

func (d *Database) Priority() int {
	return 0
}
