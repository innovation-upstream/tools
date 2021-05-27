package template_hydrator

import (
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model_frame_path"
)

type Generator interface {
	Generate(mod model_frame_path.ModelFramePath) (map[model_frame_path.ModelFrameLayerType]string, error)
}
