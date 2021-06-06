package model_frame_path

import "gitlab.innovationup.stream/innovation-upstream/tools/gen-model-frame/internal/module"

// Represent a collection of model frames and how they relate/interact with
// each other
// Can be used to generate a chain of related functions (io, logic, data)
type ModelFramePath struct {
	FunctionType module.ModelFunctionLabel `json:"type"`
	// What reference is used in this frame path
	ReferenceType ReferenceType                 `json:"referenceType"`
	Layers        []module.ModelFrameLayerLabel `json:"layers"`
}
