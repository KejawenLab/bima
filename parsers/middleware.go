package parsers

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const MIDDLEWARES_FILE = "configs/middlewares.yaml"

type Middleware struct {
	Config []string `yaml:"middlewares"`
}

func (m *Middleware) Parse(dir string) []string {
	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, MIDDLEWARES_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}

	return m.Config
}
