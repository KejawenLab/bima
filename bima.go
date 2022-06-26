package bima

import (
	"github.com/KejawenLab/bima/v3/handlers"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/models"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/utils"
	"gorm.io/gorm"
)

const (
	Version = "v3.0.4"

	HighestPriority = 255
	LowestPriority  = -255
)

type (
	Module struct {
		Handler   *handlers.Handler
		Logger    *loggers.Logger
		Cache     *utils.Cache
		Paginator *paginations.Pagination
		Request   *paginations.Request
	}

	GormModel struct {
		models.GormBase
	}

	Server struct {
		Debug    bool
		Database *gorm.DB
	}
)
