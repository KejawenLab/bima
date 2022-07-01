package bima

import (
	"github.com/KejawenLab/bima/v3/handlers"
	"github.com/KejawenLab/bima/v3/models"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/utils"
	"gorm.io/gorm"
)

const (
	Version = "v3.1.1"

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
		Debug    bool
		Database *gorm.DB
	}
)
