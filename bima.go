package bima

import (
	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/KejawenLab/bima/v2/paginations"
	"github.com/KejawenLab/bima/v2/utils"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

const VERSION_STRING = "v2.5.0"

type (
	Module struct {
		Elasticsearch *elastic.Client
		Handler       *handlers.Handler
		Logger        *handlers.Logger
		Cache         *utils.Cache
		Paginator     *paginations.Pagination
		Request       *paginations.Request
	}

	GormModel struct {
		configs.GormBase
	}

	Server struct {
		Debug    bool
		Database *gorm.DB
	}
)
