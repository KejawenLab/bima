package configs

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	MuxMiddlewares func(http.Handler) http.Handler

	Filter struct {
		Field    string
		Operator string
		Value    interface{}
	}

	User struct {
		Id    string
		Email string
		Role  int
	}

	Service struct {
		Name           string
		ConnonicalName string
		Host           string
	}

	Db struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		Driver   string
	}

	Elasticsearch struct {
		Host  string
		Port  int
		Index string
	}

	MongoDb struct {
		Host     string
		Port     int
		Database string
	}

	Amqp struct {
		Host     string
		Port     int
		User     string
		Password string
	}

	AuthHeader struct {
		Id        string
		Email     string
		Role      string
		Whitelist string
		MaxRole   int
	}

	Env struct {
		Debug            bool
		HtppPort         int
		RpcPort          int
		Version          string
		ApiVersion       string
		Service          Service
		Db               Db
		Elasticsearch    Elasticsearch
		MongoDb          MongoDb
		Amqp             Amqp
		AuthHeader       AuthHeader
		CacheLifetime    int
		User             *User
		TemplateLocation string
		RequestIDHeader  string
	}

	Template struct {
		ApiVersion            string
		PackageName           string
		Module                string
		ModuleLowercase       string
		ModulePlural          string
		ModulePluralLowercase string
		Columns               []*FieldTemplate
	}

	ModuleJson struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	ModuleTemplate struct {
		Name   string
		Fields []*FieldTemplate
	}

	FieldTemplate struct {
		Name           string
		NameUnderScore string
		ProtobufType   string
		GolangType     string
		Index          int
		IsRequired     bool
	}

	LoggerExtension struct {
		Extensions []logrus.Hook
	}

	Type struct {
		Map map[string]string
	}

	GormBase struct {
		Id        string `gorm:"type:string;primaryKey;autoIncrement:false"`
		CreatedAt sql.NullTime
		UpdatedAt sql.NullTime
		SyncedAt  sql.NullTime
		CreatedBy sql.NullString
		UpdatedBy sql.NullString
		DeletedAt gorm.DeletedAt
		DeletedBy sql.NullString
		Env       *Env `gorm:"-:all"`
	}

	MongoBase struct {
		mgm.DefaultModel `bson:",inline"`
		CreatedAt        time.Time `bson:"created_at"`
		UpdatedAt        time.Time `bson:"updated_at"`
		SyncedAt         time.Time `bson:"synced_at"`
		CreatedBy        string    `bson:"created_by"`
		UpdatedBy        string    `bson:"updated_by"`
	}
)

func (l *LoggerExtension) Register(extensions []logrus.Hook) {
	l.Extensions = extensions
}

func (t *Type) List() map[string]string {
	return t.Map
}

func (t *Type) Value(key string) string {
	if value, ok := t.Map[key]; ok {
		return value
	}

	return ""
}

func (b *GormBase) SetCreatedBy(user *User) {
	b.CreatedBy = sql.NullString{String: user.Id, Valid: true}
}

func (b *GormBase) SetUpdatedBy(user *User) {
	b.UpdatedBy = sql.NullString{String: user.Id, Valid: true}
}

func (b *GormBase) SetDeletedBy(user *User) {
	b.DeletedBy = sql.NullString{String: user.Id, Valid: true}
}

func (b *GormBase) SetCreatedAt(time time.Time) {
	b.CreatedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *GormBase) SetUpdatedAt(time time.Time) {
	b.UpdatedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *GormBase) SetSyncedAt(time time.Time) {
	b.SyncedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *GormBase) SetDeletedAt(time time.Time) {
	b.DeletedAt = gorm.DeletedAt{Time: time, Valid: true}
}

func (b *GormBase) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = uuid.New().String()

	b.SetCreatedBy(b.Env.User)
	b.SetCreatedAt(time.Now())

	return nil
}

func (b *GormBase) BeforeUpdate(tx *gorm.DB) (err error) {
	b.SetUpdatedBy(b.Env.User)
	b.SetUpdatedAt(time.Now())

	return nil
}

func (b *GormBase) BeforeDelete(tx *gorm.DB) (err error) {
	b.SetDeletedBy(b.Env.User)

	return nil
}
