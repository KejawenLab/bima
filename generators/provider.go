package generators

import (
	"fmt"
	"os"
	"strings"
)

const ModuleImport = "@modules:import"
const ModuleRegister = "@modules:register"

type Provider struct {
}

func (p *Provider) Generate(template *Template, modulePath string, packagePath string, templatePath string) {
	workDir, _ := os.Getwd()

	var path strings.Builder
	path.WriteString(workDir)
	path.WriteString("/configs/provider.go")

	file, err := os.ReadFile(path.String())
	if err != nil {
		panic(err)
	}

	contents := strings.Split(string(file), "\n")
	importIdx := 0
	moduleIdx := 0
	skipImport := true
	for k, v := range contents {
		if strings.Contains(v, ModuleImport) {
			importIdx = k
			skipImport = false
			continue
		}

		if strings.Contains(v, ModuleRegister) {
			moduleIdx = k
			break
		}
	}

	if !skipImport {
		contents[importIdx] = fmt.Sprintf(`    //%s
    %s %q`, ModuleImport, template.ModuleLowercase, fmt.Sprintf("%s/%s", template.PackageName, template.ModulePluralLowercase))
	}

	contents[moduleIdx] = fmt.Sprintf(`
    /*@module:%s*/if err := p.AddDefSlice(%s.Dic); err != nil {return err}
    //%s`, template.ModuleLowercase, template.ModuleLowercase, ModuleRegister)

	body := strings.Join(contents, "\n")

	os.WriteFile(path.String(), []byte(body), 0644)
}
