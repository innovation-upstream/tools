package analyze

import (
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"
)

type (
	//go:generate mockgen -destination=../mock/model_analyzer_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/analyze ModelAnalyzer
	ModelAnalyzer interface {
		GetModules() ([]*module.ModelFrameModule, error)
		GetDependencyModules(prevLayerModules []*module.ModelFrameModule) ([]*module.ModelFrameModule, error)
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

// TODO: use go routines to make this super fast
func (a *modelAnalyzer) GetDependencyModules(prevLayerModules []*module.ModelFrameModule) ([]*module.ModelFrameModule, error) {
	var modules []*module.ModelFrameModule

	var prevLayerDepLabels []label.ModelFrameResourceLabel
	for _, m := range prevLayerModules {
		for _, l := range m.Layers {
			for _, d := range l.Deps {
				// Prevent infinite loop when a module layer depends on sibling layers
				if d.Label.GetNamespace() != m.Name.GetNamespace() {
					prevLayerDepLabels = append(prevLayerDepLabels, d.Label)
				}
			}
		}
	}

	if len(prevLayerDepLabels) > 0 {
		currentLayerModules, err := a.ModuleLoader.LoadModules(prevLayerDepLabels)
		if err != nil {
			return modules, errors.WithStack(err)
		}

		modules = append(modules, currentLayerModules...)

		depModules, err := a.GetDependencyModules(currentLayerModules)
		if err != nil {
			return modules, errors.WithStack(err)
		}

		modules = append(modules, depModules...)
	}

	return modules, nil
}

func (a *modelAnalyzer) GetModules() ([]*module.ModelFrameModule, error) {
	var modules []*module.ModelFrameModule

	var modulesToLoad []label.ModelFrameResourceLabel
	for _, fp := range a.Model.FramePaths {
		for _, modelFrameModuleName := range fp.Layers {
			modulesToLoad = append(modulesToLoad, modelFrameModuleName)
		}
	}

	// load the modules the model directly depeonds on
	directModules, err := a.ModuleLoader.LoadModules(modulesToLoad)
	if err != nil {
		return modules, errors.WithStack(err)
	}

	// Recursive load transitive modules (modules that direct modules depend on)
	depModules, err := a.GetDependencyModules(directModules)
	if err != nil {
		return modules, errors.WithStack(err)
	}

	// concat direct modules with all dep modules
	modules = append(directModules, depModules...)

	return modules, nil
}
