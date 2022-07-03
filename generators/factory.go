package generators

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"golang.org/x/mod/modfile"
)

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
		ApiVersion string
		Driver     string
		Pluralizer *pluralize.Client
		Template   *Template
		Generators []Generator
	}
)

func (f *Factory) Generate(module ModuleTemplate) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	packagePath, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		panic(err)
	}

	workDir, _ := os.Getwd()
	packageName := f.packageName(workDir)
	moduleName := strcase.ToCamel(module.Name)
	modulePlural := f.Pluralizer.Plural(module.Name)
	modulePluralLowercase := strcase.ToDelimited(modulePlural, '_')

	var modulePath strings.Builder

	modulePath.WriteString(workDir)
	modulePath.WriteString("/")
	modulePath.WriteString(modulePluralLowercase)

	f.Template.ApiVersion = f.ApiVersion
	f.Template.PackageName = packageName
	f.Template.Module = moduleName
	f.Template.ModuleLowercase = strcase.ToDelimited(module.Name, '_')
	f.Template.ModulePlural = modulePlural
	f.Template.ModulePluralLowercase = modulePluralLowercase
	f.Template.Columns = module.Fields

	var templatePath strings.Builder

	templatePath.WriteString("templates")
	switch f.Driver {
	case "mongo":
		templatePath.WriteString("/mongo")
	default:
		templatePath.WriteString("/gorm")
	}

	os.MkdirAll(modulePath.String(), 0755)
	for _, generator := range f.Generators {
		generator.Generate(f.Template, modulePath.String(), packagePath, templatePath.String())
	}
}

func (f *Factory) packageName(workDir string) string {
	var path strings.Builder

	path.WriteString(workDir)
	path.WriteString("/go.mod")

	mod, err := os.ReadFile(path.String())
	if err != nil {
		panic(err)
	}

	return modfile.ModulePath(mod)
}
