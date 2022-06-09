package parsers

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const LISTENERS_FILE = "configs/listeners.yaml"

type Listener struct {
	Config []string `yaml:"listeners"`
}

func (l Listener) Parse(dir string) []string {
	config, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, LISTENERS_FILE))
	if err != nil {
		log.Println(err)

		return []string{}
	}

	err = yaml.Unmarshal(config, &l)
	if err != nil {
		log.Println(err)

		return []string{}
	}

	return l.Config
}
