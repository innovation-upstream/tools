package module

import (
	"github.com/pkg/errors"
)

type (
	ModuleTemplates struct {
		Templates map[ModelFunctionLabel]TemplatesForFunctionType
		Module    *ModelFrameModule
	}

	TemplatesForFunctionType struct {
		LayerTemplates  map[ModelFrameLayerLabel]TemplatesForLayer
		LayoutTemplates map[ModelFrameLayerLabel]string
	}

	TemplatesForLayer struct {
		SectionTemplates map[string]string
	}
)

func NewModuleTemplates(module *ModelFrameModule) (*ModuleTemplates, error) {
	var mod ModuleTemplates
	templates := make(map[ModelFunctionLabel]TemplatesForFunctionType)
	loader := NewModuleLoader()

	for _, f := range module.Functions {
		var funcTemplates TemplatesForFunctionType
		layerTemplates := make(map[ModelFrameLayerLabel]TemplatesForLayer)
		layoutTemplates := make(map[ModelFrameLayerLabel]string)

	layers:
		for _, l := range module.Layers {
		layerFunctions:
			for _, lf := range l.Functions {
				if lf.GetName() == f.Label.GetName() {
					break layerFunctions
				}
				continue layers
			}

			sectionTemplates := make(map[string]string)
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

func (r *ModuleTemplates) GetTemplatesForFunctionLabel(funcType ModelFunctionLabel) TemplatesForFunctionType {
	return r.Templates[ModelFunctionLabel(funcType.GetName())]
}
