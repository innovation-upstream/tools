package renderer

import (
	"github.com/tools/gen-model-frame/core/label"
	"github.com/tools/gen-model-frame/core/model"
)

type RenderedModelLayers map[label.ModelFrameResourceLabel]string

type ModelLayerRenderer interface {
	GenerateCodeLayersForFramePath(layers model.ModelLayers) (map[label.ModelFrameResourceLabel]RenderedModelLayers, error)
}
