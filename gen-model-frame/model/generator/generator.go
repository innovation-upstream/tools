package generator

import (
	"gitlab.innovationup.stream/innovation-upstream/gen-model-frame/model/model_function"
)

type Generator interface {
	Generate(mod model_function.ModelFramePath) (map[model_function.ModelFrameLayerType]string, error)
}
