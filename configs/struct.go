package configs

import "github.com/sirupsen/logrus"

type (
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
