package parsers

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const MIDDLEWARES_FILE = "configs/middlewares.yaml"

type middleware struct {
	Config []string `yaml:"middlewares"`
}

func ParseMiddleware(dir string) []string {
	config, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, MIDDLEWARES_FILE))
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
