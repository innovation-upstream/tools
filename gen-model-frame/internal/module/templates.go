package module

import "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"

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
