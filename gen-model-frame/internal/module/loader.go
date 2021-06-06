package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry"
)

//go:generate mockgen -destination=../mock/module_loader_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module ModuleLoader
type (
	ModuleLoader interface {
		LoadSectionTemplate(moduleName string, functionLabel ModelFunctionLabel, layerLabel ModelFrameLayerLabel, sectionLabel string) (string, error)
		LoadLayerLayoutTemplate(moduleName string, layerLabel ModelFrameLayerLabel) (string, error)
		LoadAllModulesFromDirectory(modulesDir string) ([]*ModelFrameModule, error)
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

func (l *moduleLoader) LoadSectionTemplate(moduleName string, functionLabel ModelFunctionLabel, layerLabel ModelFrameLayerLabel, sectionLabel string) (string, error) {
	var tmpl string

	sectionTemplatePath := fmt.Sprintf("modules/%s/templates/%s/%s/%s.gotemplate", moduleName, layerLabel, functionLabel, sectionLabel)
	by, err := ioutil.ReadFile(sectionTemplatePath)
	if err != nil {
		return tmpl, errors.WithStack(err)
	}

	tmpl = string(by)

	return tmpl, nil
}

func (l *moduleLoader) LoadLayerLayoutTemplate(moduleName string, layerLabel ModelFrameLayerLabel) (string, error) {
	var tmpl string

	sectionTemplatePath := fmt.Sprintf("modules/%s/templates/%s/layout.gotemplate", moduleName, layerLabel)
	by, err := ioutil.ReadFile(sectionTemplatePath)
	if err != nil {
		return tmpl, errors.WithStack(err)
	}

	tmpl = string(by)

	return tmpl, nil
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
			errMsg := fmt.Sprintf("failed loading %s", mJSON)
			return modules, errors.WithMessage(errors.WithStack(err), errMsg)
		}

		modules = append(modules, &module)
	}

	return modules, nil
}

func (l *moduleLoader) loadModule(moduleDir string, moduleName string) (ModelFrameModule, error) {
	var module ModelFrameModule

	moduleJSONPath := fmt.Sprintf("%s/%s/module.json", moduleDir, moduleName)
	by, err := ioutil.ReadFile(moduleJSONPath)
	if err != nil {
		return module, errors.WithStack(err)
	}

	err = json.Unmarshal(by, &module)
	if err != nil {
		errMsg := fmt.Sprintf("failed loading %s", moduleJSONPath)
		return module, errors.WithMessage(errors.WithStack(err), errMsg)
	}

	return module, nil
}
