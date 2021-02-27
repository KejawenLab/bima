package parsers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

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

	if ok, _ := regexp.MatchString(`tests$`, workDir); ok {
		workDir = path.Dir(workDir)
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
