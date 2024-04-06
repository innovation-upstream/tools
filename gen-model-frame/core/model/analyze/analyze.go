package analyze

import (
	"github.com/pkg/errors"
	"github.com/tools/gen-model-frame/core/label"
	"github.com/tools/gen-model-frame/core/model"
	"github.com/tools/gen-model-frame/core/module"
	moduleRegistryClient "github.com/tools/gen-model-frame/core/registry/module/client"
)

type (
	ModelAnalyzer interface {
		GetModules() ([]*module.ModelFrameModule, error)
		GetDependencyModules(prevLayerModules []*module.ModelFrameModule) ([]*module.ModelFrameModule, error)
	}

	modelAnalyzer struct {
		Model        model.Model
		ModuleLoader moduleRegistryClient.ModuleLoader
	}
)

func NewModelAnalyzer(m model.Model, l moduleRegistryClient.ModuleLoader) ModelAnalyzer {
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
			for _, i := range l.Implementations {
				for _, d := range i.Deps {
					// Prevent infinite loop when a module layer depends on sibling layers
					if d.LayerLabel.GetNamespace() != m.Name {
						prevLayerDepLabels = append(prevLayerDepLabels, d.LayerLabel)
					}
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
		for _, layerImplDep := range fp.Layers {
			modulesToLoad = append(modulesToLoad, layerImplDep.LayerLabel)
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
