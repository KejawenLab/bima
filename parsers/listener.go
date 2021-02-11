package parsers

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const LISTENERS_FILE = "configs/listeners.yaml"

type Listeners struct {
	Config []string `yaml:"listeners"`
}

func (l *Listeners) Parse() []string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", workDir, LISTENERS_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &l)
	if err != nil {
		panic(err)
	}

	return l.Config
}
