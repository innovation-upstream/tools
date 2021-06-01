package generator

import (
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/module"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/template_hydrator"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/transform"
)

type (
	GolangModelGenerator interface {
		Generate(m model.Model) (map[module.ModelFunctionLabel]map[string]map[module.ModelFrameLayerLabel]string, error)
	}
	golangModelGenerator struct{}
)

func NewGolangModelGenerator() GolangModelGenerator {
	return &golangModelGenerator{}
}

func (g *golangModelGenerator) Generate(m model.Model) (map[module.ModelFunctionLabel]map[string]map[module.ModelFrameLayerLabel]string, error) {
	out := make(map[module.ModelFunctionLabel]map[string]map[module.ModelFrameLayerLabel]string)

	moduleLoader := module.NewModuleLoader()
	modules, err := moduleLoader.LoadModules("modules")
	if err != nil {
		return out, errors.WithStack(err)
	}

	var templatesForModulesForModel []*module.ModuleTemplates
	for _, modu := range modules {
		for _, modModule := range m.Modules {
			if modModule == modu.Name {
				moduleTemplates, err := module.NewModuleTemplates(modu)
				if err != nil {
					return out, errors.WithStack(err)
				}

				templatesForModulesForModel = append(templatesForModulesForModel, moduleTemplates)
			}
		}
	}

	tr := transform.NewModelFramePathGoTemplateTransformer(&m)
	for _, n := range m.FramePaths {
		gen := template_hydrator.NewTemplateHydrator(tr, templatesForModulesForModel)
		res, err := gen.Generate(n)
		if err != nil {
			return out, errors.WithStack(err)
		}

		out[n.FunctionType] = res
	}

	return out, nil
}
