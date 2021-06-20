package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module/registry"
)

//go:generate mockgen -destination=../mock/module_loader_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module ModuleLoader
type (
	ModuleLoader interface {
		LoadSectionTemplate(moduleName string, functionLabel label.ModelFrameResourceLabel, layerLabel label.ModelFrameResourceLabel, sectionLabel label.ModelFrameResourceLabel) (string, error)
		LoadLayerLayoutTemplate(moduleName string, layerLabel label.ModelFrameResourceLabel) (string, error)
		LoadAllModulesFromDirectory(modulesDir string) ([]*ModelFrameModule, error)
		LoadModuleTemplates(m *ModelFrameModule) (*ModuleTemplates, error)
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

func (l *moduleLoader) LoadSectionTemplate(moduleName string, functionLabel label.ModelFrameResourceLabel, layerLabel label.ModelFrameResourceLabel, sectionLabel label.ModelFrameResourceLabel) (string, error) {
	var tmpl string

	sectionTemplatePath := fmt.Sprintf("modules/%s/templates/%s/%s/%s.gotemplate", moduleName, layerLabel.GetResourceName(), functionLabel.GetResourceName(), sectionLabel.GetResourceName())
	by, err := ioutil.ReadFile(sectionTemplatePath)
	if err != nil {
		return tmpl, errors.WithStack(err)
	}

	tmpl = string(by)

	return tmpl, nil
}

func (l *moduleLoader) LoadLayerLayoutTemplate(moduleName string, layerLabel label.ModelFrameResourceLabel) (string, error) {
	var tmpl string

	sectionTemplatePath := fmt.Sprintf("modules/%s/templates/%s/layout.gotemplate", moduleName, layerLabel.GetResourceName())
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
			errMsg := fmt.Sprintf("failed loading %s", h.GetName())
			return modules, errors.WithMessage(errors.WithStack(err), errMsg)
		}

		modules = append(modules, &module)
	}

	return modules, nil
}

func (loader *moduleLoader) LoadModuleTemplates(module *ModelFrameModule) (*ModuleTemplates, error) {
	var mod ModuleTemplates
	templates := make(map[label.ModelFrameResourceLabel]TemplatesForFunctionType)

	// TODO: iterate all loaded modules to match functions to allow modules to
	// use functions from other modules
	for _, f := range module.Functions {
		var funcTemplates TemplatesForFunctionType
		layerTemplates := make(map[label.ModelFrameResourceLabel]TemplatesForLayer)
		layoutTemplates := make(map[label.ModelFrameResourceLabel]string)

	layers:
		for _, l := range module.Layers {
		layerFunctions:
			for _, lf := range l.Functions {
				if lf.Label == f.Label {
					break layerFunctions
				}
				continue layers
			}

			sectionTemplates := make(map[label.ModelFrameResourceLabel]string)
			for _, s := range l.Sections {
				sectionTmpl, err := loader.LoadSectionTemplate(string(module.Name), f.Label, l.Label, s.Label)
				if err != nil {
					return &mod, errors.WithStack(err)
				}

				sectionTemplates[s.Label] = sectionTmpl
			}

			layerTemplates[l.Label] = TemplatesForLayer{
				SectionTemplates: sectionTemplates,
			}

			layerTmpl, err := loader.LoadLayerLayoutTemplate(string(module.Name), l.Label)
			if err != nil {
				return &mod, errors.WithStack(err)
			}

			layoutTemplates[l.Label] = layerTmpl
		}

		funcTemplates.LayoutTemplates = layoutTemplates
		funcTemplates.LayerTemplates = layerTemplates
		templates[f.Label] = funcTemplates
	}

	mod = ModuleTemplates{
		Templates: templates,
		Module:    module,
	}

	return &mod, nil
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
