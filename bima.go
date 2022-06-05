package bima

import (
	"context"

	configs "github.com/KejawenLab/bima/v2/configs"
	handlers "github.com/KejawenLab/bima/v2/handlers"
	paginations "github.com/KejawenLab/bima/v2/paginations"
	utils "github.com/KejawenLab/bima/v2/utils"
	elastic "github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

const VERSION_STRING = "v2.2.16"

type (
	Module struct {
		Context       context.Context
		Elasticsearch *elastic.Client
		Handler       *handlers.Handler
		Logger        *handlers.Logger
		Messenger     *handlers.Messenger
		Cache         *utils.Cache
		Paginator     *paginations.Pagination
		Request       *paginations.Request
	}

	GormModel struct {
		configs.GormBase
	}

	Server struct {
		Env      *configs.Env
		Database *gorm.DB
	}
)
