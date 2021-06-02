package module

import "regexp"

type (
	ModelFrameModuleName string

	ModelFrameModule struct {
		Name      ModelFrameModuleName `json:"name"`
		Functions []ModelFunction      `json:"functions"`
		Layers    []ModelFrameLayer    `json:"layers"`
	}

	ModelSection struct {
		Label string `json:"label"`
	}

	ModelFunction struct {
		Label ModelFunctionLabel `json:"label"`
	}

	ModelFrameLayer struct {
		Label     ModelFrameLayerLabel `json:"label"`
		Functions []ModelFunctionLabel `json:"functions"`
		Sections  []ModelSection       `json:"sections"`
	}
)

func (n ModelFrameModuleName) GetFileFriendlyName() string {
	moduleName := string(n)
	var reAt = regexp.MustCompile(`^@`)
	moduleName = reAt.ReplaceAllString(moduleName, "")
	var reSlash = regexp.MustCompile(`\/`)
	moduleName = reSlash.ReplaceAllString(moduleName, "_")
	var reDash = regexp.MustCompile(`-`)
	moduleName = reDash.ReplaceAllString(moduleName, "_")

	return moduleName
}
