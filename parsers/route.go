package parsers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"gopkg.in/yaml.v2"
)

const ROUTES_FILE = "configs/routes.yaml"

type Route struct {
	Config []string `yaml:"routes"`
}

func (r *Route) Parse() []string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if ok, _ := regexp.MatchString(`tests$`, workDir); ok {
		workDir = path.Dir(workDir)
	}

	config, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", workDir, ROUTES_FILE))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(config, &r)
	if err != nil {
		panic(err)
	}

	return r.Config
}
