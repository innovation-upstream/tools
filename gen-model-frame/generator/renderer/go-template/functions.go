package renderer

import (
	"text/template"

	"github.com/iancoleman/strcase"
	"innovationup.stream/tools/gen-model-frame/core/label"
)

func TxtFuncMap() template.FuncMap {
	return template.FuncMap(funcMap)
}

var funcMap = map[string]interface{}{
	"getModelNameFromLabel": getModelNameFromLabel,
	"lowerCamel":            lowerCamel,
}

func getModelNameFromLabel(name string) string {
	return label.ModelLabel(name).GetName()
}

func lowerCamel(str string) string {
	return strcase.ToLowerCamel(str)
}
