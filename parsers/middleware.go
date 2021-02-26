package parsers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"gopkg.in/yaml.v2"
)

const MIDDLEWARES_FILE = "configs/middlewares.yaml"

type Middleware struct {
	Config []string `yaml:"middlewares"`
}

func (m *Middleware) Parse() []string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if ok, _ := regexp.MatchString(`tests$`, workDir); ok {
		workDir = path.Dir(workDir)
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
