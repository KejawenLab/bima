package generators

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	configs "github.com/crowdeco/bima/configs"
	"github.com/crowdeco/bima/utils"
	"github.com/gertd/go-pluralize"
	"golang.org/x/mod/modfile"
)

const TEMPLATE_PATH = "templates"

type Factory struct {
	Env        *configs.Env
	Pluralizer *pluralize.Client
	Template   *configs.Template
	Generators []configs.Generator
	Word       *utils.Word
}

func (f *Factory) Generate(module *configs.ModuleTemplate) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	packagePath, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		panic(err)
	}

	workDir, _ := os.Getwd()
	packageName := f.GetPackageName(workDir)
	moduleName := f.Word.Camelcase(module.Name)
	modulePlural := f.Pluralizer.Plural(moduleName)
	modulePluralLowercase := f.Word.Underscore(modulePlural)
	modulePath := fmt.Sprintf("%s/%s", workDir, modulePluralLowercase)

	f.Template.ApiVersion = f.Env.ApiVersion
	f.Template.PackageName = packageName
	f.Template.Module = moduleName
	f.Template.ModuleLowercase = f.Word.Underscore(module.Name)
	f.Template.ModulePlural = modulePlural
	f.Template.ModulePluralLowercase = modulePluralLowercase
	f.Template.Columns = module.Fields

	os.MkdirAll(modulePath, 0755)
	for _, generator := range f.Generators {
		generator.Generate(f.Template, modulePath, packagePath, f.Env.TemplateLocation)
	}
}

func (f *Factory) GetDefaultTemplatePath() string {
	return TEMPLATE_PATH
}

func (f *Factory) GetPackageName(workDir string) string {
	mod, err := ioutil.ReadFile(fmt.Sprintf("%s/go.mod", workDir))
	if err != nil {
		panic(err)
	}

	return modfile.ModulePath(mod)
}
