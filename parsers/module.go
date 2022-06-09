package parsers

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const MODULES_FILE = "configs/modules.yaml"

type module struct {
	Config []string `yaml:"modules"`
}

func ParseModule(dir string) []string {
	config, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, MODULES_FILE))
	mapping := module{}
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
