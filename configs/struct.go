package configs

import "github.com/sirupsen/logrus"

type (
	User struct {
		Id    string
		Email string
		Role  int
	}

	Env struct {
		Debug                bool
		HtppPort             int
		RpcPort              int
		Version              string
		ApiVersion           string
		ServiceName          string
		ServiceCanonicalName string
		ServiceHost          string
		DbHost               string
		DbPort               int
		DbUser               string
		DbPassword           string
		DbName               string
		DbDriver             string
		ElasticsearchHost    string
		ElasticsearchPort    int
		ElasticsearchIndex   string
		MongoDbHost          string
		MongoDbPort          int
		MongoDbName          string
		AmqpHost             string
		AmqpPort             int
		AmqpUser             string
		AmqpPassword         string
		HeaderUserId         string
		HeaderUserEmail      string
		HeaderUserRole       string
		MaximumRole          int
		CacheLifetime        int
		User                 *User
		TemplateLocation     string
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
)

func (l *LoggerExtension) Register(extensions []logrus.Hook) {
	l.Extensions = extensions
}
