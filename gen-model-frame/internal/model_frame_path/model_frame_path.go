package model_frame_path

import "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"

// Represent a collection of model frames and how they relate/interact with
// each other
// Can be used to generate a chain of related functions (io, logic, data)
type ModelFramePath struct {
	Layers []label.ModelFrameResourceLabel `json:"layers"`
}
