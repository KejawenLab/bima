package bima

import (
	"github.com/KejawenLab/bima/v4/handlers"
	"github.com/KejawenLab/bima/v4/messengers"
	"github.com/KejawenLab/bima/v4/models"
	"github.com/KejawenLab/bima/v4/paginations"
	"github.com/KejawenLab/bima/v4/utils"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

const (
	Version = "v4.0.0"

	HighestPriority = 255
	LowestPriority  = -255

	Application = "application"
	Generator   = "generator"
)

type (
	Module struct {
		Debug     bool
		Handler   *handlers.Handler
		Cache     *utils.Cache
		Paginator *paginations.Pagination
	}

	GormModel struct {
		models.GormBase
	}

	Server struct {
		Debug bool
	}
)

func (s *Server) Consume(messenger *messengers.Messenger) {
}

func (s *Server) Sync(client *elastic.Client) {
}

func (s *Server) Migrate(db *gorm.DB) {
}
