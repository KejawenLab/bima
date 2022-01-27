package configs

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

const HIGEST_PRIORITY = 255

const LOWEST_PRIORITY = -255

type (
	Driver interface {
		Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB
	}

	Generator interface {
		Generate(template *Template, modulePath string, packagePath string, templatePath string)
	}

	Listener interface {
		Handle(event interface{})
		Listen() string
		Priority() int
	}

	Model interface {
		TableName() string
		SetCreatedBy(user *User)
		SetUpdatedBy(user *User)
		SetDeletedBy(user *User)
		SetCreatedAt(time time.Time)
		SetUpdatedAt(time time.Time)
		SetSyncedAt(time time.Time)
		SetDeletedAt(time time.Time)
		IsSoftDelete() bool
	}

	Module interface {
		Consume()
		Populete()
	}

	Server interface {
		RegisterGRpc(server *grpc.Server)
		GRpcHandler(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error
		RegisterAutoMigrate()
		RegisterQueueConsumer()
		RepopulateData()
	}

	Router interface {
		Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn)
		Priority() int
	}

	Parser interface {
		Parse(dir string) []string
	}

	Route interface {
		Path() string
		Method() string
		Handle(w http.ResponseWriter, r *http.Request, params map[string]string)
		SetClient(client *grpc.ClientConn)
		Middlewares() []Middleware
	}

	Middleware interface {
		Attach(request *http.Request, response http.ResponseWriter) bool
		Priority() int
	}

	Application interface {
		Run(servers []Server)
		IsBackground() bool
		Priority() int
	}
)
