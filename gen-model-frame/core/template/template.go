package template

import "innovationup.stream/tools/gen-model-frame/core/label"

type (
	ModuleTemplates struct {
		Templates map[label.ModelFrameResourceLabel]TemplatesForLayerImplementations
	}

	TemplatesForLayerImplementations struct {
		LayerImplementationTemplates map[label.ModelFrameResourceLabel]TemplatesForLayerImplementation
	}

	TemplatesForLayerImplementation struct {
		SectionTemplates map[label.ModelFrameResourceLabel]string
		LayoutTemplate   string
	}
)
