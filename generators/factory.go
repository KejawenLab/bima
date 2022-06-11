package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"golang.org/x/mod/modfile"
)

const TEMPLATE_PATH = "templates"

type (
	Generator interface {
		Generate(template *Template, modulePath string, packagePath string, templatePath string)
	}

	Template struct {
		ApiVersion            string
		PackageName           string
		Module                string
		ModuleLowercase       string
		ModulePlural          string
		ModulePluralLowercase string
		Columns               []*FieldTemplate
	}

	ModuleJson struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	ModuleTemplate struct {
		Name   string
		Fields []*FieldTemplate
	}

	FieldTemplate struct {
		Name           string
		NameUnderScore string
		ProtobufType   string
		GolangType     string
		Index          int
		IsRequired     bool
	}

	Factory struct {
		ApiVersion       string
		TemplateLocation string
		Driver           string
		Pluralizer       *pluralize.Client
		Template         *Template
		Generators       []Generator
	}
)

func (f *Factory) Generate(module *ModuleTemplate) {
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
	moduleName := strcase.ToCamel(module.Name)
	modulePlural := f.Pluralizer.Plural(module.Name)
	modulePluralLowercase := strcase.ToDelimited(modulePlural, '_')
	modulePath := fmt.Sprintf("%s/%s", workDir, modulePluralLowercase)

	f.Template.ApiVersion = f.ApiVersion
	f.Template.PackageName = packageName
	f.Template.Module = moduleName
	f.Template.ModuleLowercase = strcase.ToDelimited(module.Name, '_')
	f.Template.ModulePlural = modulePlural
	f.Template.ModulePluralLowercase = modulePluralLowercase
	f.Template.Columns = module.Fields

	templatePath := fmt.Sprintf("%s/gorm", f.TemplateLocation)
	if f.Driver == "mongo" {
		templatePath = fmt.Sprintf("%s/mongo", f.TemplateLocation)
	}

	os.MkdirAll(modulePath, 0755)
	for _, generator := range f.Generators {
		generator.Generate(f.Template, modulePath, packagePath, templatePath)
	}
}

func (f *Factory) GetDefaultTemplatePath() string {
	return TEMPLATE_PATH
}

func (f *Factory) GetPackageName(workDir string) string {
	mod, err := os.ReadFile(fmt.Sprintf("%s/go.mod", workDir))
	if err != nil {
		panic(err)
	}

	return modfile.ModulePath(mod)
}
