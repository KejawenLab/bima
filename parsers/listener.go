package parsers

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const LISTENERS_FILE = "configs/listeners.yaml"

type listener struct {
	Config []string `yaml:"listeners"`
}

func ParseListener(dir string) []string {
	config, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, LISTENERS_FILE))
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
