package parsers

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type route struct {
	Config []string `yaml:"routes"`
}

func ParseRoute(dir string) []string {
	var path strings.Builder
	path.WriteString(dir)
	path.WriteString("/")
	path.WriteString("configs/routes.yaml")

	config, err := os.ReadFile(path.String())
	mapping := route{}
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
