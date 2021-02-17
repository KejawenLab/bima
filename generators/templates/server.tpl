package {{.ModulePluralLowercase}}

import (
    "context"

	bima "github.com/crowdeco/bima"
	grpcs "{{.PackageName}}/protos/builds"
	models "{{.PackageName}}/{{.ModulePluralLowercase}}/models"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
    *bima.Server
	Module *Module
}

func (s *Server) RegisterGRpc(gs *grpc.Server) {
	grpcs.Register{{.Module}}sServer(gs, s.Module)
}

func (s *Server) GRpcHandler(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error {
	return grpcs.Register{{.Module}}sHandler(context, server, client)
}

func (s *Server) RegisterAutoMigrate() {
	if s.Env.DbAutoMigrate {
		s.Database.AutoMigrate(&models.{{.Module}}{})
	}
}

func (s *Server) RegisterQueueConsumer() {
	s.Module.Consume()
}

func (s *Server) RepopulateData() {
	if s.Env.Debug {
		s.Module.Populate()
	}
}
