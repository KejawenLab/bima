package configs

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
		MinRole   int
	}

	Env struct {
		Debug            bool
		HttpPort         int
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
)
