package {{.ModulePluralLowercase}}

import (
    "context"

	"github.com/KejawenLab/bima/v2"
	"{{.PackageName}}/protos/builds"
	"{{.PackageName}}/{{.ModulePluralLowercase}}/models"
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
	if s.Database != nil && s.Env.Debug {
		s.Database.AutoMigrate(&models.{{.Module}}{})
	}
}

func (s *Server) RegisterQueueConsumer() {
}

func (s *Server) RepopulateData() {
}
