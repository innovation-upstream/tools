package renderer

import (
	"innovationup.stream/tools/gen-model-frame/core/label"
	"innovationup.stream/tools/gen-model-frame/core/model"
)

type RenderedCodeLayers map[label.ModelFrameResourceLabel]string

type CodeLayerRenderer interface {
	GenerateCodeLayersForFramePath(mod model.ModelLayers) ([]RenderedCodeLayers, error)
}
