package model

import (
	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
)

type ModelMetadata string

type Model struct {
	Name       string                            `json:"name"`
	MetaData   map[string]string                 `json:"metadata"`
	FramePaths []model_frame_path.ModelFramePath `json:"framePaths"`
	Output     ModelOutput                       `json:"output"`
}

type ModelOutput struct {
	Directory string `json:"directory"`
}
