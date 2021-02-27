package parsers

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const UPGRADE_FILE = "configs/upgrades.yaml"

type Upgrade struct {
	Config []string `yaml:"upgrades"`
}

func (u *Upgrade) Parse(dir string) []string {
	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, UPGRADE_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &u)
	if err != nil {
		panic(err)
	}

	return u.Config
}
