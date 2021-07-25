package code_layer_generator

import (
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
)

type CodeLayers map[label.ModelFrameResourceLabel]string

//go:generate mockgen -destination=../mock/code_layer_generator_mock.go -package=mock gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/code_layer_generator CodeLayerGenerator
type CodeLayerGenerator interface {
	GenerateCodeLayersForFramePath(mod model_frame_path.ModelFramePath) ([]CodeLayers, error)
}
