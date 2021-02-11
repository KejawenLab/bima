package parsers

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const MODULES_FILE = "configs/modules.yaml"

type Module struct {
	Config []string `yaml:"modules"`
}

func (m *Module) Parse() []string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", workDir, MODULES_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}

	return m.Config
}
