package template_hydrator

import (
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/model_frame_path"
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/internal/module"
)

type Generator interface {
	Generate(mod model_frame_path.ModelFramePath) (map[string]map[module.ModelFrameLayerLabel]string, error)
}
