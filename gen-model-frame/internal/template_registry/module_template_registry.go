package template_registry

import (
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"
)

type ModuleTemplateRegistry interface {
	GetTemplatesForModelFrameLayer(layer module.ModelFrameLayerLabel) (module.TemplatesForLayer, string)
}

type moduleTemplateRegistry struct {
	ModelFramePath model_frame_path.ModelFramePath
}
