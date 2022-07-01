package configs

import (
	"context"

	"github.com/KejawenLab/bima/v3/messengers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type (
	Server interface {
		Register(server *grpc.Server)
		Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error
		Migrate(db *gorm.DB)
		Consume(messenger *messengers.Messenger)
		RepopulateData(client *elastic.Client)
	}

	User struct {
		Id    string
		Email string
		Role  int
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

	Elasticsearch struct {
		Host  string `json:"host" yaml:"host"`
		Port  int    `json:"port" yaml:"port"`
		Index string `json:"index" yaml:"index"`
	}

	Amqp struct {
		Host     string `json:"host" yaml:"host"`
		Port     int    `json:"port" yaml:"port"`
		User     string `json:"user" yaml:"user"`
		Password string `json:"password" yaml:"password"`
	}

	AuthHeader struct {
		Id        string `json:"id" yaml:"id"`
		Email     string `json:"email" yaml:"email"`
		Role      string `json:"role" yaml:"role"`
		Whitelist string `json:"whitelist" yaml:"whitelist"`
		MinRole   int    `json:"min_role" yaml:"min_role"`
	}

	Env struct {
		Debug           bool          `json:"debug" yaml:"debug"`
		Secret          string        `json:"secret" yaml:"secret"`
		HttpPort        int           `json:"http_port" yaml:"http_port"`
		RpcPort         int           `json:"rpc_port" yaml:"rpc_port"`
		ApiVersion      string        `json:"api_version" yaml:"api_version"`
		Service         Service       `json:"service" yaml:"service"`
		Db              Db            `json:"database" yaml:"database"`
		Elasticsearch   Elasticsearch `json:"elasticsearch" yaml:"elasticsearch"`
		Amqp            Amqp          `json:"queue" yaml:"queue"`
		AuthHeader      AuthHeader    `json:"auth_header" yaml:"auth_header"`
		RequestIDHeader string        `json:"request_id_header" yaml:"request_id_header"`
		CacheLifetime   int           `json:"cache_lifetime" yaml:"cache_lifetime"`
		User            *User
	}
)
