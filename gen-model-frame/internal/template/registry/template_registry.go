package registry

import "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"

type (
	TemplateRegistry interface {
		LoadTemplateForAllSections(layerLabel label.ModelFrameResourceLabel) (map[label.ModelFrameResourceLabel]string, error)
		LoadLayerLayoutTemplate(layerLabel label.ModelFrameResourceLabel) (string, error)
	}
)
