package model

import (
	"github.com/tools/gen-model-frame/core/label"
)

type (
	Model struct {
		Label      label.ModelLabel  `json:"name"`
		FramePaths []ModelLayers     `json:"framePaths"`
		Hooks      map[string]string `json:"hooks"`
		Options    ModelOptions      `json:"options"`
	}
)
