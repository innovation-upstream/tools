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

func (loader *templateLoader) LoadLayerTemplates(modelFunction label.ModelFrameResourceLabel, l ModelLayer) (TemplatesForLayer, string, error) {
	var layerLayoutTmpl string
	var templatesForLayer TemplatesForLayer

	sectionTemplates := make(map[label.ModelFrameResourceLabel]string)
	for _, s := range l.Sections {
		sectionTmpl, err := loader.registry.LoadSectionTemplate(modelFunction, l.Label, s.Label)
		if err != nil {
			return templatesForLayer, layerLayoutTmpl, errors.WithStack(err)
		}

		sectionTemplates[s.Label] = sectionTmpl
	}

	layerLayoutTmpl, err := loader.registry.LoadLayerLayoutTemplate(l.Label)
	if err != nil {
		return templatesForLayer, layerLayoutTmpl, errors.WithStack(err)
	}

	templatesForLayer = TemplatesForLayer{
		SectionTemplates: sectionTemplates,
	}

	return templatesForLayer, layerLayoutTmpl, nil
}

func (loader *templateLoader) getDependancyLayersToLoad(l ModelLayer) []ModelLayer {
	var layersToLoad []ModelLayer

	for _, d := range l.Deps {
		// find the layer that dep refs
		for _, am := range loader.modules {
			for _, al := range am.Layers {
				if al.Label == d.Label {
					layersToLoad = append(layersToLoad, al)
					depLayersToLoad := loader.getDependancyLayersToLoad(al)
					layersToLoad = append(layersToLoad, depLayersToLoad...)
				}
			}
		}
	}

	return layersToLoad
}

func (loader *templateLoader) LoadModuleTemplates(fp model_frame_path.ModelFramePath) (*ModuleTemplates, error) {
	var mod ModuleTemplates
	templates := make(map[label.ModelFrameResourceLabel]TemplatesForFunctionType)

	f := fp.Function
	var layersToLoad []ModelLayer
	for _, l := range fp.Layers {
		// get the layer struct instance for the frame path layer ref
		for _, m := range loader.modules {
			if m.Name.GetNamespace() == l.GetNamespace() {
				for _, ml := range m.Layers {
					if ml.Label == l {
						layerDepLayers := loader.getDependancyLayersToLoad(ml)
						layersToLoad = append(layersToLoad, ml)
						layersToLoad = append(layersToLoad, layerDepLayers...)
					}
				}
			}
		}
	}

	var funcTemplates TemplatesForFunctionType
	layerTemplates := make(map[label.ModelFrameResourceLabel]TemplatesForLayer)
	layoutTemplates := make(map[label.ModelFrameResourceLabel]string)

	for _, l := range layersToLoad {
		tmpltesForLayer, layerLayoutTmpl, err := loader.LoadLayerTemplates(f, l)
		if err != nil {
			return &mod, errors.WithStack(err)
		}

		layerTemplates[l.Label] = tmpltesForLayer
		layoutTemplates[l.Label] = layerLayoutTmpl
	}

	funcTemplates.LayoutTemplates = layoutTemplates
	funcTemplates.LayerTemplates = layerTemplates
	templates[f] = funcTemplates

	mod = ModuleTemplates{
		Templates: templates,
	}

	return &mod, nil
}
