package registry

import "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"

type (
	TemplateRegistry interface {
		LoadSectionTemplate(functionLabel label.ModelFrameResourceLabel, layerLabel label.ModelFrameResourceLabel, sectionLabel label.ModelFrameResourceLabel) (string, error)
		LoadLayerLayoutTemplate(layerLabel label.ModelFrameResourceLabel) (string, error)
	}
)
