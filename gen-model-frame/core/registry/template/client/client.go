package client

import (
	"github.com/pkg/errors"
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/model"
	"innovationup.stream/tools/gen-model-frame/core/module"
	host "innovationup.stream/tools/gen-model-frame/core/registry/template/host"
	"innovationup.stream/tools/gen-model-frame/core/template"
	"unknwon.dev/clog/v2"
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

func (loader *templateLoader) LoadLayerImplementationTemplates(impl module.ModelLayerImplementation) (template.TemplatesForLayerImplementation, error) {
	var layerLayoutTmpl string
	var templatesForLayer template.TemplatesForLayerImplementation

	sectionTemplates, err := loader.registry.LoadTemplateForAllSections(impl.ForLayer, impl.Label)
	if err != nil {
		return templatesForLayer, errors.WithStack(err)
	}

	layerLayoutTmpl, err = loader.registry.LoadLayerLayoutTemplate(impl.ForLayer, impl.Label)
	if err != nil {
		return templatesForLayer, errors.WithStack(err)
	}

	templatesForLayer = template.TemplatesForLayerImplementation{
		SectionTemplates: sectionTemplates,
		LayoutTemplate:   layerLayoutTmpl,
	}

	return templatesForLayer, nil
}

// TODO: prevent deuplicate layers/implementations from being loaded here
func (loader *templateLoader) getDependancyLayersToLoad(i module.ModelLayerImplementation) ([]module.ModelLayerImplementation, []module.ModelLayer) {
	var layerImplementationsToLoad []module.ModelLayerImplementation
	var looseLayerDeps []module.ModelLayer

	layerImplementationsToLoad = append(layerImplementationsToLoad, i)

top:
	for _, d := range i.Deps {
		// find the implementation for dep
		for _, am := range loader.modules {
			if am.Name == d.LayerLabel.GetNamespace() {
				for _, al := range am.Layers {
					if al.Label == d.LayerLabel {
						if d.ImplementationLabel == "" {
							looseLayerDeps = append(looseLayerDeps, al)
							continue top
						} else {
							for _, li := range al.Implementations {
								if li.Label == d.ImplementationLabel {
									depLayersToLoad, depLooseLayerDeps := loader.getDependancyLayersToLoad(li)
									for _, v := range depLayersToLoad {
										layerImplementationsToLoad = append(layerImplementationsToLoad, v)
									}
									looseLayerDeps = append(looseLayerDeps, depLooseLayerDeps...)
									continue top
								}
							}
						}
					}
				}
			}
		}
		clog.Error("No implementation found for dependancy: %s", d.LayerLabel)
	}

	return layerImplementationsToLoad, looseLayerDeps
}

// TODO: break this up
func (loader *templateLoader) LoadModuleTemplates(fp model.ModelLayers) (*template.ModuleTemplates, error) {
	var mod template.ModuleTemplates
	templates := make(map[label.ModelFrameResourceLabel]template.TemplatesForLayerImplementations)

	var allLayerImplementationsToLoad []module.ModelLayerImplementation
	var allLooseLayerDeps []module.ModelLayer
	// find the layer entrypoints to build our dep tree from
	for _, l := range fp.Layers {
		for _, m := range loader.modules {
			if m.Name == l.LayerLabel.GetNamespace() {
				for _, ml := range m.Layers {
					if ml.Label == l.LayerLabel {
						for _, i := range ml.Implementations {
							if i.Label == l.ImplementationLabel {
								depLayers, looseLayerDeps := loader.getDependancyLayersToLoad(i)
								allLayerImplementationsToLoad = append(allLayerImplementationsToLoad, depLayers...)
								allLooseLayerDeps = append(allLooseLayerDeps, looseLayerDeps...)
							}
						}
					}
				}
			}
		}
	}

	// Verfiy we loaded an implementation for every dep layer
top:
	for _, l := range allLooseLayerDeps {
		for _, m := range loader.modules {
			if m.Name == l.Label.GetNamespace() {
				for _, ml := range m.Layers {
					if ml.Label == l.Label {
						for _, li := range l.Implementations {
							for _, i := range allLayerImplementationsToLoad {
								if i.Label == li.Label {
									for _, d := range li.Deps {
										if i.Label == d.ImplementationLabel {
											continue top
										}
									}
									continue top
								}
							}
						}
					}
				}
			}
		}
		clog.Fatal("No implementation loaded for dependancy on: %s", l.Label)
	}

	// De-dupe layers
	var layerImplementationsToLoad []module.ModelLayerImplementation
all:
	for _, li := range allLayerImplementationsToLoad {
		for _, i := range layerImplementationsToLoad {
			if i.Label == li.Label && i.ForLayer == li.ForLayer {
				continue all
			}
		}
		layerImplementationsToLoad = append(layerImplementationsToLoad, li)
	}

	clog.Trace("Loading templates for %+v", layerImplementationsToLoad)
	// Actually load the templates
	for _, i := range layerImplementationsToLoad {
		tmpltesForLayerImpl, err := loader.LoadLayerImplementationTemplates(i)
		if err != nil {
			return &mod, errors.WithStack(err)
		}

		templates[i.ForLayer] = template.TemplatesForLayerImplementations{
			LayerImplementationTemplates: map[label.ModelFrameResourceLabel]template.TemplatesForLayerImplementation{
				i.Label: tmpltesForLayerImpl,
			},
		}
	}

	mod = template.ModuleTemplates{
		Templates: templates,
	}

	return &mod, nil
}
