package parsers

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const LISTENERS_FILE = "configs/listeners.yaml"

type Listeners struct {
	Config []string `yaml:"listeners"`
}

func (l *Listeners) Parse(dir string) []string {
	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, LISTENERS_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &l)
	if err != nil {
		panic(err)
	}

	return l.Config
}
