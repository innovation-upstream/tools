package module

import (
	"github.com/pkg/errors"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/template/registry"
)

type (
	TemplateLoader interface {
		LoadModuleTemplates(fp model_frame_path.ModelFramePath) (*ModuleTemplates, error)
	}

	templateLoader struct {
		modules  []*ModelFrameModule
		registry registry.TemplateRegistry
	}
)

func NewTemplateLoader(modules []*ModelFrameModule, registry registry.TemplateRegistry) TemplateLoader {
	return &templateLoader{
		modules,
		registry,
	}
}

func (loader *templateLoader) LoadLayerTemplates(l ModelLayer) (TemplatesForLayer, error) {
	var layerLayoutTmpl string
	var templatesForLayer TemplatesForLayer

	sectionTemplates, err := loader.registry.LoadTemplateForAllSections(l.Label)
	if err != nil {
		return templatesForLayer, errors.WithStack(err)
	}

	layerLayoutTmpl, err = loader.registry.LoadLayerLayoutTemplate(l.Label)
	if err != nil {
		return templatesForLayer, errors.WithStack(err)
	}

	templatesForLayer = TemplatesForLayer{
		SectionTemplates: sectionTemplates,
		LayoutTemplate:   layerLayoutTmpl,
	}

	return templatesForLayer, nil
}

func (loader *templateLoader) getDependancyLayersToLoad(l ModelLayer) []ModelLayer {
	var layersToLoad []ModelLayer

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

func (loader *templateLoader) LoadModuleTemplates(fp model_frame_path.ModelFramePath) (*ModuleTemplates, error) {
	var mod ModuleTemplates
	var templates []TemplatesForLayers

	var allLayersToLoad []ModelLayer
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
	var layersToLoad []ModelLayer
all:
	for _, al := range allLayersToLoad {
		for _, l := range layersToLoad {
			if l.Label == al.Label {
				continue all
			}
		}
		layersToLoad = append(layersToLoad, al)
	}

	var templatesForLayers TemplatesForLayers
	layerTemplates := make(map[label.ModelFrameResourceLabel]TemplatesForLayer)

	for _, l := range layersToLoad {
		tmpltesForLayer, err := loader.LoadLayerTemplates(l)
		if err != nil {
			return &mod, errors.WithStack(err)
		}

		layerTemplates[l.Label] = tmpltesForLayer
	}

	templatesForLayers.LayerTemplates = layerTemplates
	templates = append(templates, templatesForLayers)

	mod = ModuleTemplates{
		Templates: templates,
	}

	return &mod, nil
}
