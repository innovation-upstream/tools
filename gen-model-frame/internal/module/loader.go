package module

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry"
)

//go:generate mockgen -destination=../mock/module_loader_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module ModuleLoader
type (
	ModuleLoader interface {
		LoadAllModulesFromDirectory(modulesDir string) ([]*ModelFrameModule, error)
		LoadModules(l []label.ModelFrameResourceLabel) ([]*ModelFrameModule, error)
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

func (l *moduleLoader) LoadAllModulesFromDirectory(modulesDir string) ([]*ModelFrameModule, error) {
	var modules []*ModelFrameModule

	headers, err := l.Registry.QueryAllModuleHeaders()
	if err != nil {
		return modules, errors.WithStack(err)
	}

	for _, h := range headers {
		var module ModelFrameModule
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

func (loader *moduleLoader) LoadModules(lbls []label.ModelFrameResourceLabel) ([]*ModelFrameModule, error) {
	var modules []*ModelFrameModule

	headers, err := loader.Registry.QueryModuleHeaders(lbls)
	if err != nil {
		return modules, errors.WithStack(err)
	}

	for _, h := range headers {
		var module ModelFrameModule
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
