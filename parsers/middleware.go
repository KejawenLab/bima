package parsers

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type middleware struct {
	Config []string `yaml:"middlewares"`
}

func ParseMiddleware(dir string) []string {
	var path strings.Builder
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
