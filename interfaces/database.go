package interfaces

import (
	"fmt"

	configs "github.com/crowdeco/bima/v2/configs"
	"github.com/fatih/color"
)

type Database struct {
}

func (d *Database) Run(servers []configs.Server) {
	color.New(color.FgCyan, color.Bold).Printf("✓ ")
	fmt.Println("Serving DB Auto Migration Juices...")

	for _, server := range servers {
		go server.RegisterAutoMigrate()
	}
}

func (d *Database) IsBackground() bool {
	return true
}

func (d *Database) Priority() int {
	return configs.HIGEST_PRIORITY + 1
}
