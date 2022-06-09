package parsers

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const ROUTES_FILE = "configs/routes.yaml"

type Route struct {
	Config []string `yaml:"routes"`
}

func (r Route) Parse(dir string) []string {
	config, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, ROUTES_FILE))
	if err != nil {
		log.Println(err)

		return []string{}
	}

	err = yaml.Unmarshal(config, &r)
	if err != nil {
		log.Println(err)

		return []string{}
	}

	return r.Config
}
