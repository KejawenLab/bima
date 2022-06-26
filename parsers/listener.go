package parsers

import (
	"bytes"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type listener struct {
	Config []string `yaml:"listeners"`
}

func ParseListener(dir string) []string {
	var path bytes.Buffer
	path.WriteString(dir)
	path.WriteString("/")
	path.WriteString("configs/listeners.yaml")

	config, err := os.ReadFile(path.String())
	mapping := listener{}
	if err != nil {
		log.Println(err)

		return []string{}
	}

	err = yaml.Unmarshal(config, &mapping)
	if err != nil {
		log.Println(err)

		return []string{}
	}

	return mapping.Config
}
