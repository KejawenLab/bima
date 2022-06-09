package dics

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/KejawenLab/bima/v2"
	"github.com/KejawenLab/bima/v2/configs"
	"github.com/KejawenLab/bima/v2/configs/drivers"
	"github.com/KejawenLab/bima/v2/events"
	"github.com/KejawenLab/bima/v2/generators"
	"github.com/KejawenLab/bima/v2/handlers"
	"github.com/KejawenLab/bima/v2/interfaces"
	"github.com/KejawenLab/bima/v2/listeners/creates"
	"github.com/KejawenLab/bima/v2/listeners/deletes"
	filters "github.com/KejawenLab/bima/v2/listeners/paginations"
	"github.com/KejawenLab/bima/v2/listeners/updates"
	"github.com/KejawenLab/bima/v2/middlewares"
	paginations "github.com/KejawenLab/bima/v2/paginations"
	"github.com/KejawenLab/bima/v2/paginations/adapter"
	"github.com/KejawenLab/bima/v2/parsers"
	"github.com/KejawenLab/bima/v2/repositories"
	"github.com/KejawenLab/bima/v2/routers"
	"github.com/KejawenLab/bima/v2/routes"
	"github.com/KejawenLab/bima/v2/utils"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kamva/mgm/v3"
	"github.com/olivere/elastic/v7"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var Container = []dingo.Def{
	{
		Name:  "bima:config:parser:listener",
		Build: (*parsers.Listener)(nil),
	},
	{
		Name:  "bima:config:parser:logger",
		Build: (*parsers.Logger)(nil),
	},
	{
		Name:  "bima:config:parser:middleware",
		Build: (*parsers.Middleware)(nil),
	},
	{
		Name:  "bima:config:parser:module",
		Build: (*parsers.Module)(nil),
	},
	{
		Name:  "bima:config:parser:route",
		Build: (*parsers.Route)(nil),
	},
	{
		Name:  "bima:config:user",
		Build: (*configs.User)(nil),
	},
	{
		Name: "bima:config:type",
		Build: func() (*configs.Type, error) {
			return &configs.Type{
				Map: map[string]string{
					"double":   "float64",
					"float":    "float32",
					"int32":    "int32",
					"int64":    "int64",
					"uint32":   "uint32",
					"uint64":   "uint64",
					"sint32":   "int32",
					"sint64":   "int64",
					"fixed32":  "uint32",
					"fixed64":  "uint64",
					"sfixed32": "int32",
					"sfixed64": "int64",
					"bool":     "bool",
					"string":   "string",
					"bytes":    "[]byte",
				},
			}, nil
		},
	},
	{
		Name:  "bima:config:template",
		Build: (*configs.Template)(nil),
	},
	{
		Name:  "bima:template:module",
		Build: (*configs.ModuleTemplate)(nil),
	},
	{
		Name:  "bima:template:field",
		Build: (*configs.FieldTemplate)(nil),
	},
	{
		Name: "bima:config:env",
		Build: func() (*configs.Env, error) {
			env := configs.Env{}

			env.Version = os.Getenv("APP_VERSION")
			env.ApiVersion = os.Getenv("API_VERSION")
			env.RequestIDHeader = os.Getenv("REQUEST_ID_HEADER")
			env.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
			env.HttpPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
			env.RpcPort, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))

			if env.RequestIDHeader == "" {
				env.RequestIDHeader = "X-Request-Id"
			}

			sName := os.Getenv("APP_NAME")
			env.Service = configs.Service{
				Name:           sName,
				ConnonicalName: utils.Underscore(sName),
				Host:           os.Getenv("APP_HOST"),
			}

			dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
			env.Db = configs.Db{
				Host:     os.Getenv("DB_HOST"),
				Port:     dbPort,
				User:     os.Getenv("DB_USER"),
				Password: os.Getenv("DB_PASSWORD"),
				Name:     os.Getenv("DB_NAME"),
				Driver:   os.Getenv("DB_DRIVER"),
			}

			err := mgm.SetDefaultConfig(nil, env.Db.Name, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", env.Db.User, env.Db.Password, env.Db.Host, env.Db.Port)).SetMonitor(&event.CommandMonitor{
				Started: func(_ context.Context, evt *event.CommandStartedEvent) {
					log.Print(evt.Command)
				},
			}))
			if err != nil {
				err = mgm.SetDefaultConfig(nil, env.Db.Name, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", env.Db.Host)).SetMonitor(&event.CommandMonitor{
					Started: func(_ context.Context, evt *event.CommandStartedEvent) {
						log.Print(evt.Command)
					},
				}))
			}

			if err != nil {
				return nil, err
			}

			esPort, _ := strconv.Atoi(os.Getenv("ELASTICSEARCH_PORT"))
			env.Elasticsearch = configs.Elasticsearch{
				Host:  os.Getenv("ELASTICSEARCH_HOST"),
				Port:  esPort,
				Index: env.Db.Name,
			}

			mgdbPort, _ := strconv.Atoi(os.Getenv("MONGODB_PORT"))
			env.MongoDb = configs.MongoDb{
				Host:     os.Getenv("MONGODB_HOST"),
				Port:     mgdbPort,
				Database: "data_logs",
			}

			amqpPort, _ := strconv.Atoi(os.Getenv("AMQP_PORT"))
			env.Amqp = configs.Amqp{
				Host:     os.Getenv("AMQP_HOST"),
				Port:     amqpPort,
				User:     os.Getenv("AMQP_USER"),
				Password: os.Getenv("AMQP_PASSWORD"),
			}

			minRole, _ := strconv.Atoi(os.Getenv("AUTH_HEADER_MIN_ROLE"))
			env.AuthHeader = configs.AuthHeader{
				Id:        os.Getenv("AUTH_HEADER_ID"),
				Email:     os.Getenv("AUTH_HEADER_EMAIL"),
				Role:      os.Getenv("AUTH_HEADER_ROLE"),
				Whitelist: os.Getenv("AUTH_HEADER_WHITELIST"),
				MinRole:   minRole,
			}

			env.CacheLifetime, _ = strconv.Atoi(os.Getenv("CACHE_LIFETIME"))
			env.TemplateLocation = generators.TEMPLATE_PATH

			return &env, nil
		},
	},
	{
		Name: "bima:module:generator",
		Build: func(
			dic configs.Generator,
			model configs.Generator,
			module configs.Generator,
			proto configs.Generator,
			provider configs.Generator,
			server configs.Generator,
			validation configs.Generator,
			swagger configs.Generator,
			env *configs.Env,
			pluralizer *pluralize.Client,
			template *configs.Template,
		) (*generators.Factory, error) {
			return &generators.Factory{
				ApiVersion:       env.ApiVersion,
				Driver:           env.Db.Driver,
				TemplateLocation: env.TemplateLocation,
				Pluralizer:       pluralizer,
				Template:         template,
				Generators: []configs.Generator{
					dic,
					model,
					module,
					proto,
					provider,
					server,
					validation,
					swagger,
				},
			}, nil
		},
		Params: dingo.Params{
			"0":  dingo.Service("bima:generator:dic"),
			"1":  dingo.Service("bima:generator:model"),
			"2":  dingo.Service("bima:generator:module"),
			"3":  dingo.Service("bima:generator:proto"),
			"4":  dingo.Service("bima:generator:provider"),
			"5":  dingo.Service("bima:generator:server"),
			"6":  dingo.Service("bima:generator:validation"),
			"7":  dingo.Service("bima:generator:swagger"),
			"8":  dingo.Service("bima:config:env"),
			"9":  dingo.Service("bima:util:pluralizer"),
			"10": dingo.Service("bima:config:template"),
		},
	},
	{
		Name:  "bima:generator:dic",
		Build: (*generators.Dic)(nil),
	},
	{
		Name:  "bima:generator:model",
		Build: (*generators.Model)(nil),
	},
	{
		Name:  "bima:generator:module",
		Build: (*generators.Module)(nil),
		Params: dingo.Params{
			"Config": dingo.Service("bima:config:parser:module"),
		},
	},
	{
		Name:  "bima:generator:proto",
		Build: (*generators.Proto)(nil),
	},
	{
		Name:  "bima:generator:provider",
		Build: (*generators.Provider)(nil),
	},
	{
		Name:  "bima:generator:server",
		Build: (*generators.Server)(nil),
	},
	{
		Name:  "bima:generator:validation",
		Build: (*generators.Validation)(nil),
	},
	{
		Name:  "bima:generator:swagger",
		Build: (*generators.Swagger)(nil),
	},
	{
		Name:  "bima:database:driver:mysql",
		Build: (*drivers.Mysql)(nil),
	},
	{
		Name:  "bima:database:driver:postgresql",
		Build: (*drivers.PostgreSql)(nil),
	},
	{
		Name: "bima:application",
		Build: func(
			database configs.Application,
			elasticsearch configs.Application,
			grpc configs.Application,
			queue configs.Application,
			rest configs.Application,
		) (*interfaces.Application, error) {
			return &interfaces.Application{
				Applications: []configs.Application{database, elasticsearch, grpc, queue, rest},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:interface:database"),
			"1": dingo.Service("bima:interface:elasticsearch"),
			"2": dingo.Service("bima:interface:grpc"),
			"3": dingo.Service("bima:interface:queue"),
			"4": dingo.Service("bima:interface:rest"),
		},
	},
	{
		Name:  "bima:event:dispatcher",
		Build: (*events.Dispatcher)(nil),
	},
	{
		Name: "bima:handler:middleware",
		Build: func(dipatcher *events.Dispatcher, logger *handlers.Logger) (*handlers.Middleware, error) {
			middleware := handlers.Middleware{
				Dispatcher: dipatcher,
				Logger:     logger,
			}
			middleware.Add(&middlewares.Header{})
			middleware.Add(&middlewares.Recovery{Logger: logger})

			return &middleware, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:event:dispatcher"),
			"1": dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name:  "bima:logger:extension",
		Build: (*configs.LoggerExtension)(nil),
	},
	{
		Name: "bima:connection:database",
		Build: func(env *configs.Env, mysql configs.Driver, postgresql configs.Driver) (*gorm.DB, error) {
			util := color.New(color.FgCyan, color.Bold)
			var db configs.Driver

			util.Printf("✓ ")
			fmt.Printf("Database configured using '%s' driver...\n", env.Db.Driver)
			time.Sleep(100 * time.Millisecond)

			switch env.Db.Driver {
			case "mysql":
				db = mysql
			case "postgresql":
				db = postgresql
			case "mongo":
				return nil, nil
			default:
				return nil, errors.New("Unknown database driver")
			}

			return db.Connect(
				env.Db.Host,
				env.Db.Port,
				env.Db.User,
				env.Db.Password,
				env.Db.Name,
				env.Debug,
			), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
			"1": dingo.Service("bima:database:driver:mysql"),
			"2": dingo.Service("bima:database:driver:postgresql"),
		},
	},
	{
		Name: "bima:connection:elasticsearch",
		Build: func(env *configs.Env) (*elastic.Client, error) {
			client, err := elastic.NewClient(
				elastic.SetURL(fmt.Sprintf("%s:%d", env.Elasticsearch.Host, env.Elasticsearch.Port)),
				elastic.SetSniff(false),
				elastic.SetHealthcheck(false),
				elastic.SetGzip(true),
			)

			if err != nil {
				return nil, err
			}

			color.New(color.FgCyan, color.Bold).Printf("✓ ")
			fmt.Println("Elasticsearch configured...")
			time.Sleep(100 * time.Millisecond)

			return client, nil
		},
	},
	{
		Name: "bima:listener:create:elasticsearch",
		Build: func(env *configs.Env, context context.Context, client *elastic.Client) (*creates.Elasticsearch, error) {
			return &creates.Elasticsearch{
				Service:       env.Service,
				Context:       context,
				Elasticsearch: client,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
			"1": dingo.Service("bima:context:background"),
			"2": dingo.Service("bima:connection:elasticsearch"),
		},
	},
	{
		Name: "bima:listener:update:elasticsearch",
		Build: func(env *configs.Env, context context.Context, client *elastic.Client) (*updates.Elasticsearch, error) {
			return &updates.Elasticsearch{
				Service:       env.Service,
				Context:       context,
				Elasticsearch: client,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
			"1": dingo.Service("bima:context:background"),
			"2": dingo.Service("bima:connection:elasticsearch"),
		},
	},
	{
		Name: "bima:listener:delete:elasticsearch",
		Build: func(env *configs.Env, context context.Context, client *elastic.Client) (*deletes.Elasticsearch, error) {
			return &deletes.Elasticsearch{
				Service:       env.Service,
				Context:       context,
				Elasticsearch: client,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
			"1": dingo.Service("bima:context:background"),
			"2": dingo.Service("bima:connection:elasticsearch"),
		},
	},
	{
		Name:  "bima:listener:filter:elasticsearch",
		Build: (*filters.ElasticsearchFilter)(nil),
	},
	{
		Name:  "bima:listener:filter:gorm",
		Build: (*filters.GormFilter)(nil),
	},
	{
		Name:  "bima:listener:filter:mongo",
		Build: (*filters.MongoDbFilter)(nil),
	},
	{
		Name:  "bima:interface:database",
		Build: (*interfaces.Database)(nil),
	},
	{
		Name:  "bima:interface:elasticsearch",
		Build: (*interfaces.Elasticsearch)(nil),
	},
	{
		Name: "bima:interface:grpc",
		Build: func(env *configs.Env, gRpc *grpc.Server, logger *handlers.Logger) (*interfaces.GRpc, error) {
			return &interfaces.GRpc{
				GRpcPort: env.RpcPort,
				GRpc:     gRpc,
				Logger:   logger,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
			"1": dingo.Service("bima:grpc:server"),
			"2": dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name:  "bima:interface:queue",
		Build: (*interfaces.Queue)(nil),
	},
	{
		Name: "bima:interface:rest",
		Build: func(
			env *configs.Env,
			middleware *handlers.Middleware,
			router *handlers.Router,
			server *http.ServeMux,
			context context.Context,
			logger *handlers.Logger,
		) (*interfaces.Rest, error) {
			return &interfaces.Rest{
				GRpcPort:   env.RpcPort,
				HttpPort:   env.HttpPort,
				Middleware: middleware,
				Router:     router,
				Server:     server,
				Context:    context,
				Logger:     logger,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
			"1": dingo.Service("bima:handler:middleware"),
			"2": dingo.Service("bima:handler:router"),
			"3": dingo.Service("bima:http:mux"),
			"4": dingo.Service("bima:context:background"),
			"5": dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name: "bima:handler:logger",
		Build: func(
			env *configs.Env,
			logger *logrus.Logger,
			extension *configs.LoggerExtension,
		) (*handlers.Logger, error) {
			logger.SetFormatter(&logrus.JSONFormatter{})
			for _, e := range extension.Extensions {
				logger.AddHook(e)
			}

			return &handlers.Logger{
				Verbose: env.Debug,
				Service: env.Service,
				Logger:  logger,
			}, nil
		},
	},
	{
		Name:  "bima:handler:messager",
		Build: (*handlers.Messenger)(nil),
		Params: dingo.Params{
			"Logger":    dingo.Service("bima:handler:logger"),
			"Publisher": dingo.Service("bima:message:publisher"),
			"Consumer":  dingo.Service("bima:message:consumer"),
		},
	},
	{
		Name: "bima:handler:router",
		Build: func(gateway configs.Router, mux configs.Router) (*handlers.Router, error) {
			return &handlers.Router{
				Routers: []configs.Router{
					gateway,
					mux,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:router:gateway"),
			"1": dingo.Service("bima:router:mux"),
		},
	},
	{
		Name:  "bima:middleware:auth",
		Build: (*middlewares.Auth)(nil),
		Params: dingo.Params{
			"Env": dingo.Service("bima:config:env"),
		},
	},
	{
		Name: "bima:middleware:requestid",
		Build: func(env *configs.Env, logger *handlers.Logger) (*middlewares.RequestID, error) {
			return &middlewares.RequestID{
				Logger:          logger,
				RequestIDHeader: env.RequestIDHeader,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
			"1": dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name: "bima:router:mux",
		Build: func(
			apiDoc configs.Route,
			apiDocRedirection configs.Route,
			health configs.Route,
		) (*routers.MuxRouter, error) {
			routers := routers.MuxRouter{}
			routers.Register([]configs.Route{apiDoc, apiDocRedirection, health})

			return &routers, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:routes:api-doc"),
			"1": dingo.Service("bima:routes:api-doc-redirect"),
			"2": dingo.Service("bima:routes:health"),
		},
	},
	{
		Name:  "bima:router:gateway",
		Build: (*routers.GRpcGateway)(nil),
	},
	{
		Name: "bima:routes:api-doc",
		Build: func(env *configs.Env) (*routes.ApiDoc, error) {
			return &routes.ApiDoc{
				Debug: env.Debug,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
		},
	},
	{
		Name:  "bima:routes:api-doc-redirect",
		Build: (*routes.ApiDocRedirect)(nil),
	},
	{
		Name:  "bima:routes:health",
		Build: (*routes.Health)(nil),
		Params: dingo.Params{
			"Logger": dingo.Service("bima:handler:logger"),
		},
	},
	{
		Name: "bima:http:mux",
		Build: func() (*http.ServeMux, error) {
			return http.NewServeMux(), nil
		},
	},
	{
		Name: "bima:grpc:server",
		Build: func() (*grpc.Server, error) {
			return grpc.NewServer(
				grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
					grpc_recovery.StreamServerInterceptor(),
				)),
				grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
					grpc_recovery.UnaryServerInterceptor(),
				)),
			), nil
		},
	},
	{
		Name: "bima:log:logger",
		Build: func() (*logrus.Logger, error) {
			color.New(color.FgCyan, color.Bold).Printf("✓ ")
			fmt.Println("Logger configured...")
			time.Sleep(100 * time.Millisecond)

			return logrus.New(), nil
		},
	},
	{
		Name: "bima:context:background",
		Build: func() (context.Context, error) {
			return context.Background(), nil
		},
	},
	{
		Name: "bima:message:config",
		Build: func(env *configs.Env) (amqp.Config, error) {
			color.New(color.FgCyan, color.Bold).Printf("✓ ")
			fmt.Println("Pub/Sub configured...")
			time.Sleep(100 * time.Millisecond)

			return amqp.NewDurableQueueConfig(fmt.Sprintf("amqp://%s:%s@%s:%d/", env.Amqp.User, env.Amqp.Password, env.Amqp.Host, env.Amqp.Port)), nil
		},
	},
	{
		Name: "bima:message:publisher",
		Build: func(env *configs.Env, config amqp.Config) (*amqp.Publisher, error) {
			publisher, err := amqp.NewPublisher(config, watermill.NewStdLogger(env.Debug, env.Debug))
			if err != nil {
				return nil, err
			}

			return publisher, nil
		},
	},
	{
		Name: "bima:message:consumer",
		Build: func(env *configs.Env, config amqp.Config) (*amqp.Subscriber, error) {
			consumer, err := amqp.NewSubscriber(config, watermill.NewStdLogger(env.Debug, env.Debug))
			if err != nil {
				return nil, err
			}

			return consumer, nil
		},
	},
	{
		Name:  "bima:pagination:paginator",
		Build: (*paginations.Pagination)(nil),
	},
	{
		Name:  "bima:pagination:adapter:gorm",
		Build: (*adapter.GormAdapter)(nil),
		Params: dingo.Params{
			"Env":        dingo.Service("bima:config:env"),
			"Database":   dingo.Service("bima:connection:database"),
			"Dispatcher": dingo.Service("bima:event:dispatcher"),
		},
	},
	{
		Name:  "bima:pagination:adapter:elasticsearch",
		Build: (*adapter.ElasticsearchAdapter)(nil),
		Params: dingo.Params{
			"Env":        dingo.Service("bima:config:env"),
			"Client":     dingo.Service("bima:connection:elasticsearch"),
			"Dispatcher": dingo.Service("bima:event:dispatcher"),
		},
	},
	{
		Name:  "bima:pagination:adapter:mongo",
		Build: (*adapter.MongodbAdapter)(nil),
		Params: dingo.Params{
			"Env":        dingo.Service("bima:config:env"),
			"Dispatcher": dingo.Service("bima:event:dispatcher"),
		},
	},
	{
		Name:  "bima:pagination:request",
		Build: (*paginations.Request)(nil),
	},
	{
		Name:  "bima:service:repository:gorm",
		Build: (*repositories.GormRepository)(nil),
		Params: dingo.Params{
			"Database": dingo.Service("bima:connection:database"),
		},
	},
	{
		Name:  "bima:service:repository:mongo",
		Build: (*repositories.MongoRepository)(nil),
	},
	{
		Name: "bima:cache:memory",
		Build: func(env *configs.Env) (*utils.Cache, error) {
			return utils.NewCache(time.Duration(env.CacheLifetime) * time.Second), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
		},
	},
	{
		Name: "bima:util:pluralizer",
		Build: func() (*pluralize.Client, error) {
			return pluralize.NewClient(), nil
		},
	},
	{
		Name:  "bima:module",
		Build: (*bima.Module)(nil),
		Params: dingo.Params{
			"Context":       dingo.Service("bima:context:background"),
			"Elasticsearch": dingo.Service("bima:connection:elasticsearch"),
			"Handler":       dingo.Service("bima:handler:handler"),
			"Logger":        dingo.Service("bima:handler:logger"),
			"Messenger":     dingo.Service("bima:handler:messager"),
			"Cache":         dingo.Service("bima:cache:memory"),
			"Paginator":     dingo.Service("bima:pagination:paginator"),
		},
	},
	{
		Name: "bima:server",
		Build: func(env *configs.Env, db *gorm.DB) (*bima.Server, error) {
			return &bima.Server{
				Debug:    env.Debug,
				Database: db,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
			"1": dingo.Service("bima:connection:database"),
		},
	},
	{
		Name: "bima:model",
		Build: func(env *configs.Env) (*bima.GormModel, error) {
			return &bima.GormModel{
				GormBase: configs.GormBase{
					Env: env,
				},
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("bima:config:env"),
		},
	},
}
