package model_frame_path

import "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/label"

// Represent a collection of model frames and how they relate/interact with
// each other
// Can be used to generate a chain of related functions (io, logic, data)
type ModelFramePath struct {
	FunctionType label.ModelFrameResourceLabel `json:"type"`
	// What reference is used in this frame path
	ReferenceType ReferenceType                   `json:"referenceType"`
	Layers        []label.ModelFrameResourceLabel `json:"layers"`
}
