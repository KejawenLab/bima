package parsers

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const ROUTES_FILE = "configs/routes.yaml"

type route struct {
	Config []string `yaml:"routes"`
}

func ParseRoute(dir string) []string {
	config, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, ROUTES_FILE))
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
