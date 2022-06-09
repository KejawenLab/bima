package parsers

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const MIDDLEWARES_FILE = "configs/middlewares.yaml"

type Middleware struct {
	Config []string `yaml:"middlewares"`
}

func (m Middleware) Parse(dir string) []string {
	config, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, MIDDLEWARES_FILE))
	if err != nil {
		log.Println(err)

		return []string{}
	}

	err = yaml.Unmarshal(config, &m)
	if err != nil {
		log.Println(err)

		return []string{}
	}

	return m.Config
}
