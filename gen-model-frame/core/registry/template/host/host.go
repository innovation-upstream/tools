package host

import "github.com/tools/gen-model-frame/core/label"

type (
	TemplateRegistry interface {
		LoadTemplateForAllSections(layerLabel label.ModelFrameResourceLabel, implementationLabel label.ModelFrameResourceLabel) (map[label.ModelFrameResourceLabel]string, error)
		LoadLayerLayoutTemplate(layerLabel label.ModelFrameResourceLabel, implementationLabel label.ModelFrameResourceLabel) (string, error)
	}
)
