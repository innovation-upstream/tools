package template

import "innovationup.stream/tools/gen-model-frame/core/label"

type (
	ModuleTemplates struct {
		Templates []TemplatesForLayers
	}

	TemplatesForLayers struct {
		LayerTemplates map[label.ModelFrameResourceLabel]TemplatesForLayer
	}

	TemplatesForLayer struct {
		SectionTemplates map[label.ModelFrameResourceLabel]string
		LayoutTemplate   string
	}
)
