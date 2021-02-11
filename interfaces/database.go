package interfaces

import (
	"log"

	configs "github.com/crowdeco/bima/configs"
)

type Database struct {
}

func (d *Database) Run(servers []configs.Server) {
	log.Printf("Starting DB Auto Migration")

	for _, server := range servers {
		server.RegisterAutoMigrate()
	}
}

func (d *Database) IsBackground() bool {
	return true
}

func (d *Database) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
