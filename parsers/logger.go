package parsers

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const LOGGERS_FILE = "configs/loggers.yaml"

type logger struct {
	Config []string `yaml:"loggers"`
}

func ParseLogger(dir string) []string {
	var path strings.Builder
	path.WriteString(dir)
	path.WriteString("/")
	path.WriteString("configs/loggers.yaml")

	config, err := os.ReadFile(path.String())
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
