package configs

import (
	"context"

	"github.com/KejawenLab/bima/v4/messengers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var Database *gorm.DB

type (
	Server interface {
		Register(server *grpc.Server)
		Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error
		Migrate(db *gorm.DB)
		Consume(messenger *messengers.Messenger)
		Sync(client *elastic.Client)
	}

	Service struct {
		Name           string `json:"name" yaml:"name"`
		ConnonicalName string
	}

	Db struct {
		Host     string `json:"host" yaml:"host"`
		Port     int    `json:"port" yaml:"port"`
		User     string `json:"user" yaml:"user"`
		Password string `json:"password" yaml:"password"`
		Name     string `json:"name" yaml:"name"`
		Driver   string `json:"driver" yaml:"driver"`
	}

	Env struct {
		Debug           bool    `json:"debug" yaml:"debug"`
		Secret          string  `json:"secret" yaml:"secret"`
		HttpPort        int     `json:"http_port" yaml:"http_port"`
		RpcPort         int     `json:"rpc_port" yaml:"rpc_port"`
		ApiVersion      string  `json:"api_version" yaml:"api_version"`
		Service         Service `json:"service" yaml:"service"`
		Db              Db      `json:"database" yaml:"database"`
		RequestIDHeader string  `json:"request_id_header" yaml:"request_id_header"`
		CacheLifetime   int     `json:"cache_lifetime" yaml:"cache_lifetime"`
		User            string
	}
)
