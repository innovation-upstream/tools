package model

import (
	"strings"

	"gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/model_frame_path"
)

type (
	Model struct {
		Label      ModelLabel                        `json:"name"`
		FramePaths []model_frame_path.ModelFramePath `json:"framePaths"`
		Options    ModelOptions                      `json:"options"`
	}

	ModelLabel string
)

func (n ModelLabel) GetNamespace() string {
	return strings.Split(string(n), "/")[0]
}

func (n ModelLabel) GetName() string {
	s := strings.Split(string(n), "/")
	if len(s) > 1 {
		return s[1]
	}

	return string(n)
}
