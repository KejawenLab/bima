package parsers

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const MIDDLEWARES_FILE = "configs/modules.yaml"

type Middleware struct {
	Config []string `yaml:"middlewares"`
}

func (m *Middleware) Parse() []string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", workDir, MIDDLEWARES_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}

	return m.Config
}
