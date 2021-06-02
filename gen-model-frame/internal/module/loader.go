package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type (
	//go:generate mockgen -destination=../mock/module_loader_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module ModuleLoader
	ModuleLoader interface {
		LoadSectionTemplate(moduleName string, functionLabel ModelFunctionLabel, layerLabel ModelFrameLayerLabel, sectionLabel string) (string, error)
		LoadLayerLayoutTemplate(moduleName string, layerLabel ModelFrameLayerLabel) (string, error)
		LoadModules(modulesDir string) ([]*ModelFrameModule, error)
	}
	moduleLoader struct{}
)

func NewModuleLoader() ModuleLoader {
	return &moduleLoader{}
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

func (l *moduleLoader) LoadModules(modulesDir string) ([]*ModelFrameModule, error) {
	var modules []*ModelFrameModule

	f, err := os.Open(modulesDir)
	if err != nil {
		return modules, errors.WithStack(err)
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return modules, errors.WithStack(err)
	}

	for _, file := range files {
		if file.IsDir() {
			moduleJSONFilePath := fmt.Sprintf("%s/%s/module.json", modulesDir, file.Name())
			// if its a namespaced module dir, load nested modules
			if _, err := os.Stat(moduleJSONFilePath); errors.Is(err, os.ErrNotExist) {
				nsModuleDir := fmt.Sprintf("%s/%s", modulesDir, file.Name())
				nsModules, err := l.LoadModules(nsModuleDir)
				if err != nil {
					return modules, errors.WithStack(err)
				}

				modules = append(modules, nsModules...)
				continue
			}

			module, err := l.loadModule(modulesDir, file.Name())
			if err != nil {
				return modules, errors.WithStack(err)
			}

			modules = append(modules, &module)
		}
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
