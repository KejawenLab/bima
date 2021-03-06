package parsers

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const ROUTES_FILE = "configs/routes.yaml"

type Route struct {
	Config []string `yaml:"routes"`
}

func (r Route) Parse(dir string) []string {
	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, ROUTES_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &r)
	if err != nil {
		panic(err)
	}

	return r.Config
}
