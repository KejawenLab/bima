package parsers

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const LOGGERS_FILE = "configs/loggers.yaml"

type logger struct {
	Config []string `yaml:"loggers"`
}

func ParseLogger(dir string) []string {
	config, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, LOGGERS_FILE))
	mapping := logger{}
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
