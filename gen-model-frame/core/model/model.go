package model

import (
	"innovationup.stream/tools/gen-model-frame/core/label"
)

type (
	Model struct {
		Label      label.ModelLabel `json:"name"`
		FramePaths []ModelLayers    `json:"framePaths"`
		Options    ModelOptions     `json:"options"`
	}
)
