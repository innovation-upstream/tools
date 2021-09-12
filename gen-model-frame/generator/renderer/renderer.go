package renderer

import (
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/model"
)

type RenderedModelLayers map[label.ModelFrameResourceLabel]string

type ModelLayerRenderer interface {
	GenerateCodeLayersForFramePath(layers model.ModelLayers) (map[label.ModelFrameResourceLabel]RenderedModelLayers, error)
}
