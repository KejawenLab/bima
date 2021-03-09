package bima

import (
	"context"

	configs "github.com/crowdeco/bima/v2/configs"
	handlers "github.com/crowdeco/bima/v2/handlers"
	paginations "github.com/crowdeco/bima/v2/paginations"
	utils "github.com/crowdeco/bima/v2/utils"
	elastic "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const VERSION_STRING = "v2.0.4"

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

	Model struct {
		configs.Base
	}

	Server struct {
		Env      *configs.Env
		Database *gorm.DB
	}

	Plugin interface {
		GetRoutes() []configs.Route
		GetLoggers() []logrus.Hook
		GetMiddlewares() []configs.Middleware
		GetListeners() []configs.Listener
		GetServers() []configs.Server
		GetUpgrades() []configs.Upgrade
		GetVersion() string
	}
)
