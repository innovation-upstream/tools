package registry

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/module"
	registry "innovationup.stream/tools/gen-model-frame/core/registry/module/host"
)

type (
	ModuleLoader interface {
		LoadAllModulesFromDirectory(modulesDir string) ([]*module.ModelFrameModule, error)
		LoadModules(l []label.ModelFrameResourceLabel) ([]*module.ModelFrameModule, error)
	}

	moduleLoader struct {
		Registry registry.ModuleRegistry
	}
)

func NewModuleLoader(registry registry.ModuleRegistry) ModuleLoader {
	return &moduleLoader{
		Registry: registry,
	}
}

func (l *moduleLoader) LoadAllModulesFromDirectory(modulesDir string) ([]*module.ModelFrameModule, error) {
	var modules []*module.ModelFrameModule

	headers, err := l.Registry.QueryAllModuleHeaders()
	if err != nil {
		return modules, errors.WithStack(err)
	}

	for _, h := range headers {
		var module module.ModelFrameModule
		mJSON, err := h.GetJSON()
		if err != nil {
			return modules, errors.WithStack(err)
		}

		err = json.Unmarshal([]byte(mJSON), &module)
		if err != nil {
			errMsg := fmt.Sprintf("failed loading %s", h.GetName())
			return modules, errors.WithMessage(errors.WithStack(err), errMsg)
		}

		modules = append(modules, &module)
	}

	return modules, nil
}

func (loader *moduleLoader) LoadModules(lbls []label.ModelFrameResourceLabel) ([]*module.ModelFrameModule, error) {
	var modules []*module.ModelFrameModule

	headers, err := loader.Registry.QueryModuleHeaders(lbls)
	if err != nil {
		return modules, errors.WithStack(err)
	}

	for _, h := range headers {
		var module module.ModelFrameModule
		mJSON, err := h.GetJSON()
		if err != nil {
			return modules, errors.WithStack(err)
		}

		err = json.Unmarshal([]byte(mJSON), &module)
		if err != nil {
			errMsg := fmt.Sprintf("failed loading %s", h.GetName())
			return modules, errors.WithMessage(errors.WithStack(err), errMsg)
		}

		module = module.FullyQualifyLabels()

		modules = append(modules, &module)
	}

	return modules, nil
}
