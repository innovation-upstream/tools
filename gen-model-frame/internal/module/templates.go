package module

import "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"

type (
	ModuleTemplates struct {
		Templates map[label.ModelFrameResourceLabel]TemplatesForFunctionType
	}

	TemplatesForFunctionType struct {
		LayerTemplates  map[label.ModelFrameResourceLabel]TemplatesForLayer
		LayoutTemplates map[label.ModelFrameResourceLabel]string
	}

	TemplatesForLayer struct {
		SectionTemplates map[label.ModelFrameResourceLabel]string
	}
)

func (r *ModuleTemplates) GetTemplatesForFunctionLabel(funcType label.ModelFrameResourceLabel) TemplatesForFunctionType {
	return r.Templates[funcType]
}
