package bima

import (
	"github.com/KejawenLab/bima/v3/handlers"
	"github.com/KejawenLab/bima/v3/messengers"
	"github.com/KejawenLab/bima/v3/models"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/utils"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

const (
	Version = "v3.2.0"

	HighestPriority = 255
	LowestPriority  = -255
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
