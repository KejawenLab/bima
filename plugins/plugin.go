package plugins

import (
	"os"
	"path/filepath"
	"plugin"

	bima "github.com/crowdeco/bima"
	configs "github.com/crowdeco/bima/v2/configs"
	"github.com/sirupsen/logrus"
)

const PLUGIN_NAME = "BimaPluginName"

const PLUGIN_INIT = "BimaPlugin"

type Plugin struct {
	plugins []bima.Plugin
}

func (p *Plugin) Scan(path string) {
	plugins := []*plugin.Plugin{}
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			return nil
		}

		if filepath.Ext(path) != ".so" {
			return nil
		}

		p, err := plugin.Open(path)
		if !info.IsDir() && err == nil {
			plugins = append(plugins, p)
		}

		return nil
	})

	p.register(plugins)
}

func (p *Plugin) GetRoutes() []configs.Route {
	var routes []configs.Route
	for _, x := range p.plugins {
		routes = append(routes, x.GetRoutes()...)
	}

	return routes
}

func (p *Plugin) GetMiddlewares() []configs.Middleware {
	var middlewares []configs.Middleware
	for _, x := range p.plugins {
		middlewares = append(middlewares, x.GetMiddlewares()...)
	}

	return middlewares
}

func (p *Plugin) GetLoggers() []logrus.Hook {
	var loggers []logrus.Hook
	for _, x := range p.plugins {
		loggers = append(loggers, x.GetLoggers()...)
	}

	return loggers
}

func (p *Plugin) GetListeners() []configs.Listener {
	var listeners []configs.Listener
	for _, x := range p.plugins {
		listeners = append(listeners, x.GetListeners()...)
	}

	return listeners
}

func (p *Plugin) GetServers() []configs.Server {
	var servers []configs.Server
	for _, x := range p.plugins {
		servers = append(servers, x.GetServers()...)
	}

	return servers
}

func (p *Plugin) GetUpgrades() []configs.Upgrade {
	var upgrades []configs.Upgrade
	for _, x := range p.plugins {
		upgrades = append(upgrades, x.GetUpgrades()...)
	}

	return upgrades
}

func (p *Plugin) List() []bima.Plugin {
	return p.plugins
}

func (p *Plugin) register(plugins []*plugin.Plugin) {
	var list []bima.Plugin
	for _, plugin := range plugins {
		_, err := plugin.Lookup(PLUGIN_NAME)
		if err != nil {
			continue
		}

		sign, err := plugin.Lookup(PLUGIN_INIT)
		if err != nil {
			continue
		}

		valid, ok := sign.(bima.Plugin)
		if !ok {
			continue
		}

		list = append(list, valid)
	}

	p.plugins = list
}
