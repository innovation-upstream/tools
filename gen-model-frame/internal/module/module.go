package module

type (
	ModelFrameModule struct {
		Name      string            `json:"name"`
		Functions []ModelFunction   `json:"functions"`
		Layers    []ModelFrameLayer `json:"layers"`
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
