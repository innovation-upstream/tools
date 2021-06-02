package analyze

import (
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"
)

type (
	//go:generate mockgen -destination=../mock/model_analyzer_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/analyze ModelAnalyzer
	ModelAnalyzer interface {
		GetModuleTemplates() ([]*module.ModuleTemplates, error)
	}

	modelAnalyzer struct {
		Model        model.Model
		ModuleLoader module.ModuleLoader
	}
)

func NewModelAnalyzer(m model.Model, l module.ModuleLoader) ModelAnalyzer {
	return &modelAnalyzer{
		Model:        m,
		ModuleLoader: l,
	}
}

func (a *modelAnalyzer) GetModuleTemplates() ([]*module.ModuleTemplates, error) {
	var templatesForModules []*module.ModuleTemplates
	// TODO: only load the modules we are going to use
	modules, err := a.ModuleLoader.LoadModules("modules")
	if err != nil {
		return templatesForModules, errors.WithStack(err)
	}

	for _, modu := range modules {
		for _, modModule := range a.Model.Modules {
			if modModule == string(modu.Name) {
				moduleTemplates, err := module.NewModuleTemplates(modu)
				if err != nil {
					return templatesForModules, errors.WithStack(err)
				}

				templatesForModules = append(templatesForModules, moduleTemplates)
			}
		}
	}

	return templatesForModules, nil
}
