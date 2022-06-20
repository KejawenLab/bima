package configs

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type (
	Server interface {
		RegisterGRpc(server *grpc.Server)
		GRpcHandler(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error
		RegisterAutoMigrate()
		RegisterQueueConsumer()
		RepopulateData()
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
)
