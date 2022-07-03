package parsers

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type listener struct {
	Config []string `yaml:"listeners"`
}

func ParseListener(dir string) []string {
	var path strings.Builder
	path.WriteString(dir)
	path.WriteString("/")
	path.WriteString("configs/listeners.yaml")

	config, err := os.ReadFile(path.String())
	mapping := listener{}
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
