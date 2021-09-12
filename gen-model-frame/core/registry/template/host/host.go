package host

import "innovationup.stream/tools/gen-model-frame/core/label"

type (
	TemplateRegistry interface {
		LoadTemplateForAllSections(layerLabel label.ModelFrameResourceLabel, implementationLabel label.ModelFrameResourceLabel) (map[label.ModelFrameResourceLabel]string, error)
		LoadLayerLayoutTemplate(layerLabel label.ModelFrameResourceLabel, implementationLabel label.ModelFrameResourceLabel) (string, error)
	}
)
