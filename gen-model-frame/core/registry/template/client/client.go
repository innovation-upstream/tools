package client

import (
	"github.com/pkg/errors"
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/model"
	"innovationup.stream/tools/gen-model-frame/core/module"
	host "innovationup.stream/tools/gen-model-frame/core/registry/template/host"
	"innovationup.stream/tools/gen-model-frame/core/template"
)

type (
	TemplateLoader interface {
		LoadModuleTemplates(fp model.ModelLayers) (*template.ModuleTemplates, error)
	}

	templateLoader struct {
		modules  []*module.ModelFrameModule
		registry host.TemplateRegistry
	}
)

func NewTemplateLoader(modules []*module.ModelFrameModule, registry host.TemplateRegistry) TemplateLoader {
	return &templateLoader{
		modules,
		registry,
	}
}

func (loader *templateLoader) LoadLayerTemplates(l module.ModelLayer) (template.TemplatesForLayer, error) {
	var layerLayoutTmpl string
	var templatesForLayer template.TemplatesForLayer

	sectionTemplates, err := loader.registry.LoadTemplateForAllSections(l.Label)
	if err != nil {
		return templatesForLayer, errors.WithStack(err)
	}

	layerLayoutTmpl, err = loader.registry.LoadLayerLayoutTemplate(l.Label)
	if err != nil {
		return templatesForLayer, errors.WithStack(err)
	}

	templatesForLayer = template.TemplatesForLayer{
		SectionTemplates: sectionTemplates,
		LayoutTemplate:   layerLayoutTmpl,
	}

	return templatesForLayer, nil
}

func (loader *templateLoader) getDependancyLayersToLoad(l module.ModelLayer) []module.ModelLayer {
	var layersToLoad []module.ModelLayer

	layersToLoad = append(layersToLoad, l)

	for _, d := range l.Deps {
		// find the layer for dep
		for _, am := range loader.modules {
			if am.Name.GetNamespace() == d.Label.GetNamespace() {
				for _, al := range am.Layers {
					if al.Label == d.Label {
						depLayersToLoad := loader.getDependancyLayersToLoad(al)
						layersToLoad = append(layersToLoad, depLayersToLoad...)
					}
				}
			}
		}
	}

	return layersToLoad
}

func (loader *templateLoader) LoadModuleTemplates(fp model.ModelLayers) (*template.ModuleTemplates, error) {
	var mod template.ModuleTemplates
	var templates []template.TemplatesForLayers

	var allLayersToLoad []module.ModelLayer
	// find the layer entrypoints to build our dep tree from
	for _, l := range fp.Layers {
		for _, m := range loader.modules {
			if m.Name.GetNamespace() == l.GetNamespace() {
				for _, ml := range m.Layers {
					if ml.Label == l {
						depLayers := loader.getDependancyLayersToLoad(ml)
						allLayersToLoad = append(allLayersToLoad, depLayers...)
					}
				}
			}
		}
	}

	// De-dupe layers
	var layersToLoad []module.ModelLayer
all:
	for _, al := range allLayersToLoad {
		for _, l := range layersToLoad {
			if l.Label == al.Label {
				continue all
			}
		}
		layersToLoad = append(layersToLoad, al)
	}

	var templatesForLayers template.TemplatesForLayers
	layerTemplates := make(map[label.ModelFrameResourceLabel]template.TemplatesForLayer)

	for _, l := range layersToLoad {
		tmpltesForLayer, err := loader.LoadLayerTemplates(l)
		if err != nil {
			return &mod, errors.WithStack(err)
		}

		layerTemplates[l.Label] = tmpltesForLayer
	}

	templatesForLayers.LayerTemplates = layerTemplates
	templates = append(templates, templatesForLayers)

	mod = template.ModuleTemplates{
		Templates: templates,
	}

	return &mod, nil
}
