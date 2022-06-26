package parsers

import (
	"bytes"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type middleware struct {
	Config []string `yaml:"middlewares"`
}

func ParseMiddleware(dir string) []string {
	var path bytes.Buffer
	path.WriteString(dir)
	path.WriteString("/")
	path.WriteString("configs/middlewares.yaml")

	config, err := os.ReadFile(path.String())
	mapping := middleware{}
	if err != nil {
		log.Println(err)

		return []string{}
	}

	err = yaml.Unmarshal(config, &mapping)
	if err != nil {
		log.Println(err)

		return []string{}
	}

	return mapping.Config
}
