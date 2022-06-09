package parsers

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const LOGGERS_FILE = "configs/loggers.yaml"

type Logger struct {
	Config []string `yaml:"loggers"`
}

func (l Logger) Parse(dir string) []string {
	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, LOGGERS_FILE))
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
