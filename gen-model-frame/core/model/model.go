package model

import (
	"strings"
)

type (
	Model struct {
		Label      ModelLabel    `json:"name"`
		FramePaths []ModelLayers `json:"framePaths"`
		Options    ModelOptions  `json:"options"`
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
