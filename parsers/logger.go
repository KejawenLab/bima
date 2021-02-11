package parsers

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const LOGGERS_FILE = "configs/loggers.yaml"

type Logger struct {
	Config []string `yaml:"loggers"`
}

func (l *Logger) Parse() []string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", workDir, LOGGERS_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &l)
	if err != nil {
		panic(err)
	}

	return l.Config
}
