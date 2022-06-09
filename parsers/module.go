package parsers

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const MODULES_FILE = "configs/modules.yaml"

type Module struct {
	Config []string `yaml:"modules"`
}

func (m Module) Parse(dir string) []string {
	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, MODULES_FILE))
	if err != nil {
		log.Println(err)

		return []string{}
	}

	err = yaml.Unmarshal(config, &m)
	if err != nil {
		log.Println(err)

		return []string{}
	}

	return m.Config
}
